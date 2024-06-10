#!/bin/bash

# Exit script if a statement returns a non-true return value.
set -o errexit
# Use the error status of the first failure, rather than that of the last item in a pipeline.
set -o pipefail
# Exit script if you try to use an uninitialized variable.
set -o nounset

# Configuration variables
CPU_WORKERS=${CPU_WORKERS:-4}
CPU_TIMEOUT=${CPU_TIMEOUT:-60s}
MEMORY_SIZE=${MEMORY_SIZE:-1024}  # in MB
MEMORY_ITERATIONS=${MEMORY_ITERATIONS:-1}
IO_NAME=${IO_NAME:-randwrite}
IO_ENGINE=${IO_ENGINE:-libaio}
IO_DEPTH=${IO_DEPTH:-4}
IO_RW=${IO_RW:-randwrite}
IO_BS=${IO_BS:-4k}
IO_SIZE=${IO_SIZE:-1G}
IO_JOBS=${IO_JOBS:-4}
IO_RUNTIME=${IO_RUNTIME:-60}

# Number of iterations of the stresstest
ITERATIONS=${1:-1}

run_cpu_stress_test() {
    echo "Running CPU stress test with ${CPU_WORKERS} workers for ${CPU_TIMEOUT}..."
    stress --cpu "${CPU_WORKERS}" --timeout "${CPU_TIMEOUT}"
}

run_memory_stress_test() {
    echo "Running memory stress test with ${MEMORY_SIZE}MB for ${MEMORY_ITERATIONS} iterations..."
    memtester "${MEMORY_SIZE}" "${MEMORY_ITERATIONS}"
}

run_io_stress_test() {
    echo "Running I/O stress test with ${IO_JOBS} jobs, ${IO_BS} block size for ${IO_RUNTIME}s..."
    fio --name="${IO_NAME}" --ioengine="${IO_ENGINE}" --iodepth="${IO_DEPTH}" --rw="${IO_RW}" --bs="${IO_BS}" --direct=1 --size="${IO_SIZE}" --numjobs="${IO_JOBS}" --runtime="${IO_RUNTIME}" --group_reporting
}

# Run stress tests
for (( i=1; i<=ITERATIONS; i++ ))
do
    echo "Iteration ${i} of ${ITERATIONS}"
    run_cpu_stress_test
    run_memory_stress_test
    run_io_stress_test
done

echo "Stress testing completed."
