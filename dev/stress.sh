#!/bin/bash

# Configuration variables
VERSION="1.0"
DURATION=300        # Duration to run the tests (in seconds)
CPU_LOAD_CORES=2    # Number of cores to stress
MEMORY_LOAD=256M    # Amount of memory to stress
FILE_SIZE=1G        # File size for I/O tests

# Function to display version
display_version() {
    echo "Stress Custom Carbonaut Script Version $VERSION"
    exit 0
}

# Check for --version argument and exit
for arg in "$@"; do
    case $arg in
        --version)
            display_version
            ;;
    esac
done

# Create a temporary directory for sysbench
TEST_DIR=$(mktemp -d)

echo "Starting load generation for $DURATION seconds..."

# Starting CPU stress
echo "Stressing CPU on $CPU_LOAD_CORES cores..."
stress --cpu $CPU_LOAD_CORES --timeout $DURATION &

# Starting Memory stress
echo "Stressing Memory with $MEMORY_LOAD..."
stress --vm 1 --vm-bytes $MEMORY_LOAD --timeout $DURATION &

# Preparing file for I/O test
echo "Preparing file of $FILE_SIZE for I/O tests..."
sysbench fileio --file-total-size=$FILE_SIZE prepare

# Starting I/O stress
echo "Starting I/O stress..."
sysbench fileio --file-total-size=$FILE_SIZE --file-test-mode=rndrw --time=$DURATION --max-requests=0 run

# Clean up
echo "Cleaning up..."
sysbench fileio --file-total-size=$FILE_SIZE cleanup
rm -rf $TEST_DIR

echo "Load generation completed."
