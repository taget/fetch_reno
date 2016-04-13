#!/bin/bash
# Author: Eli Qiao <qiaoliyong@gmail.com>

PROJS=$(cat  proj.txt)
LOG_FILE="$(date +%x)"

if [ -f "${LOG_FILE}" ]; then
    rm -rf "${LOG_FILE}"
fi

for item in ${PROJS[@]}; do
    ./fetchreno "${item}" >> "${LOG_FILE}"
done
