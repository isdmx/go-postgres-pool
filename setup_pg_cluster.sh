#!/bin/bash -e

PG_USER="test"
PG_PASSWORD="secret"
PG_DB="test"

# Variables
NETWORK_NAME="pg-cluster"
MASTER_NAME="pg-master"
STANDBY_NAME="pg-standby"
MASTER_PORT=5432
STANDBY_PORT=5433
PG_IMAGE=postgres:alpine

# Create Docker network
echo "Creating Docker network..."
docker network create $NETWORK_NAME

# Create Docker volumes
echo "Creating Docker volumes..."
docker volume create pg-master-data
docker volume create pg-standby-data

# Start Master Node
echo "Starting Master Node..."
docker run -d \
  --name $MASTER_NAME \
  --network $NETWORK_NAME \
  -e POSTGRES_USER=$PG_USER \
  -e POSTGRES_PASSWORD=$PG_PASSWORD \
  -e POSTGRES_DB=$PG_DB \
  -v pg-master-data:/var/lib/postgresql/data \
  -p $MASTER_PORT:5432 \
  --pull always \
 ${PG_IMAGE} 

# Wait for Master to initialize
echo "Waiting for Master to initialize..."
sleep 5

# Configure Master for Replication
echo "Configuring Master for replication..."
docker exec -it $MASTER_NAME bash -c "
  echo 'wal_level = replica' >> /var/lib/postgresql/data/postgresql.conf;
  echo 'max_wal_senders = 10' >> /var/lib/postgresql/data/postgresql.conf;
  echo 'wal_keep_size = 1GB' >> /var/lib/postgresql/data/postgresql.conf;
  echo 'host replication $PG_USER 0.0.0.0/0 md5' >> /var/lib/postgresql/data/pg_hba.conf;
"

# Restart Master to apply configuration changes
echo "Restarting Master to apply configuration changes..."
docker restart $MASTER_NAME

# Wait for Master to restart
echo "Waiting for Master to restart..."
sleep 5

# Start Standby Node
echo "Starting Standby Node..."
docker run -d \
  --name $STANDBY_NAME \
  --network $NETWORK_NAME \
  -e POSTGRES_USER=$PG_USER \
  -e POSTGRES_PASSWORD=$PG_PASSWORD \
  -v pg-standby-data:/var/lib/postgresql/data \
  -p $STANDBY_PORT:5432 \
  --pull always \
 ${PG_IMAGE} 

# Wait for Standby to initialize
echo "Waiting for Standby to initialize..."
sleep 5

# Configure Standby Node
echo "Configuring Standby Node..."
docker exec -it $STANDBY_NAME bash -c "
  rm -rf /var/lib/postgresql/data/*;
  export PGPASSWORD='$PG_PASSWORD';  # Set the password for pg_basebackup
  pg_basebackup -h $MASTER_NAME -U $PG_USER -D /var/lib/postgresql/data -P -R -X stream -C -S pgstandby1;
  echo 'hot_standby = on' >> /var/lib/postgresql/data/postgresql.conf;
  touch /var/lib/postgresql/data/standby.signal;
"

# Restart Standby to apply configuration changes
echo "Restarting Standby to apply configuration changes..."
docker restart $STANDBY_NAME

echo "PostgreSQL cluster setup complete!"
echo "Master is running on port $MASTER_PORT"
echo "Standby is running on port $STANDBY_PORT"
