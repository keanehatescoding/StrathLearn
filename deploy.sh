#!/bin/bash
set -e

echo "Loading Docker images..."
docker load < /root/runner.tar
docker load < /root/app.tar

echo "Creating required volumes..."
docker volume create --name=code-runner-data

echo "Stopping old containers..."
docker-compose down || true

echo "Starting new containers..."
docker-compose up -d

echo "Cleaning up..."
rm -f /root/runner.tar /root/app.tar

echo "Deployment complete!"