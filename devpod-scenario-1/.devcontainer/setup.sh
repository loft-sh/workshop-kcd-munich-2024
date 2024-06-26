#!/usr/bin/env bash

# Check if we can run unit tests successfully 
make test-unit

CLUSTER_NAME="dev"

# Create a development kind cluster if we don't have one yet 
for cluster in $(kind get clusters); do
    if [ "$cluster" == $CLUSTER_NAME ]; then
        echo "Cluster '$CLUSTER_NAME' found! Nothing left to do"
        break
    fi
done

if [ "$cluster" != $CLUSTER_NAME ]; then
    echo "Cluster '$CLUSTER_NAME' not found, creating..."
    kind create cluster --name $CLUSTER_NAME
fi

# Validate we can access the cluster
echo "Listing all namespaces of cluster '$CLUSTER_NAME'"
kubectl get namespaces
