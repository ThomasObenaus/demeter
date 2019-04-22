# Configuration

## Global

### Config-File

|         |                                                                      |
| ------- | -------------------------------------------------------------------- |
| name    | config-file                                                          |
| usage   | Specifies the full path and name of the configuration file for sokar |
| type    | string                                                               |
| default | ""                                                                   |
| flag    | --config-file                                                        |
| env     | -                                                                    |

### Dry-Run

|         |                                                                                                                                      |
| ------- | ------------------------------------------------------------------------------------------------------------------------------------ |
| name    | dry-run                                                                                                                              |
| usage   | If true, then sokar won't execute the planned scaling action. Only scaling actions triggered via ScaleBy end-point will be executed. |
| type    | bool                                                                                                                                 |
| default | false                                                                                                                                |
| flag    | --dry-run                                                                                                                            |
| env     | SK_DRY_RUN                                                                                                                           |

### Port

|         |                                |
| ------- | ------------------------------ |
| name    | port                           |
| usage   | Port where sokar is listening. |
| type    | uint                           |
| default | 11000                          |
| flag    | --port                         |
| env     | SK_PORT                        |

## Nomad

### Server-Address

|         |                                            |
| ------- | ------------------------------------------ |
| name    | server-address                             |
| usage   | Specifies the address of the nomad server. |
| type    | string                                     |
| default | ""                                         |
| flag    | --nomad.server-address                     |
| env     | SK_NOMAD_SERVER_ADDRESS                    |

## Job

### Name

|         |                                   |
| ------- | --------------------------------- |
| name    | name                              |
| usage   | The name of the job to be scaled. |
| type    | string                            |
| default | ""                                |
| flag    | --job.name                        |
| env     | SK_JOB_NAME                       |

### Min

|         |                               |
| ------- | ----------------------------- |
| name    | min                           |
| usage   | The minimum scale of the job. |
| type    | uint                          |
| default | 1                             |
| flag    | --job.min                     |
| env     | SK_JOB_MIN                    |

### Max

|         |                               |
| ------- | ----------------------------- |
| name    | max                           |
| usage   | The maximum scale of the job. |
| type    | uint                          |
| default | 10                            |
| flag    | --job.max                     |
| env     | SK_JOB_MAX                    |

## CapacityPlanner

### Down Scale Cooldown

|         |                                                          |
| ------- | -------------------------------------------------------- |
| name    | down-scale-cooldown                                      |
| usage   | The time sokar waits between downscaling actions at min. |
| type    | duration                                                 |
| default | 20s                                                      |
| flag    | --cap.down-scale-cooldown                                |
| env     | SK_CAP_DOWN_SCALE_COOLDOWN                               |

### Up Scale Cooldown

|         |                                                        |
| ------- | ------------------------------------------------------ |
| name    | up-scale-cooldown                                      |
| usage   | The time sokar waits between upscaling actions at min. |
| type    | duration                                               |
| default | 20s                                                    |
| flag    | --cap.up-scale-cooldown                                |
| env     | SK_CAP_UP_SCALE_COOLDOWN                               |

## Logging

### Structured

|         |                                |
| ------- | ------------------------------ |
| name    | logging.structured             |
| usage   | Use structured logging or not. |
| type    | bool                           |
| default | false                          |
| flag    | --logging.structured           |
| env     | SK_LOGGING_STRUCTURED          |

### Unix Timestamp

|         |                                                    |
| ------- | -------------------------------------------------- |
| name    | logging.unix-ts                                    |
| usage   | Use Unix-Timestamp representation for log entries. |
| type    | bool                                               |
| default | false                                              |
| flag    | --logging.unix-ts                                  |
| env     | SK_LOGGING_UNIX_TS                                 |

### NoColor

|         |                                                 |
| ------- | ----------------------------------------------- |
| name    | logging.no-color                                |
| usage   | If true colors in log out-put will be disabled. |
| type    | bool                                            |
| default | false                                           |
| flag    | --logging.no-color                              |
| env     | SK_LOGGING_NO_COLOR                             |

## ScaleAlertAggregator

### No Alert Damping

|         |                                                                                          |
| ------- | ---------------------------------------------------------------------------------------- |
| name    | no-alert-damping                                                                         |
| usage   | Damping used in case there are currently no alerts firing (neither down- nor upscaling). |
| type    | float                                                                                    |
| default | 1.0                                                                                      |
| flag    | --saa.no-alert-damping                                                                   |
| env     | SK_SAA_NO_ALERT_DAMPING                                                                  |

### Up Scale Threshold

|         |                                  |
| ------- | -------------------------------- |
| name    | up-thresh                        |
| usage   | Threshold for a upscaling event. |
| type    | float                            |
| default | 10.0                             |
| flag    | --saa.up-thresh                  |
| env     | SK_SAA_UP_THRESH                 |

### Down Scale Threshold

|         |                                    |
| ------- | ---------------------------------- |
| name    | down-thresh                        |
| usage   | Threshold for a downscaling event. |
| type    | float                              |
| default | 10.0                               |
| flag    | --saa.down-thresh                  |
| env     | SK_SAA_DOWN_THRESH                 |

### Evaluation Cycle

|         |                                                                                                 |
| ------- | ----------------------------------------------------------------------------------------------- |
| name    | eval-cycle                                                                                      |
| usage   | Cycle/ frequency the ScaleAlertAggregator evaluates the weights of the currently firing alerts. |
| type    | Duration                                                                                        |
| default | 1s                                                                                              |
| flag    | --saa.eval-cycle                                                                                |
| env     | SK_SAA_EVAL_CYCLE                                                                               |

### Evaluation Period Factor

|         |                                                                                                                                  |
| ------- | -------------------------------------------------------------------------------------------------------------------------------- |
| name    | eval-period-factor                                                                                                               |
| usage   | EvaluationPeriodFactor is used to calculate the evaluation period (evaluationPeriod = evaluationCycle * evaluationPeriodFactor). |
| type    | uint                                                                                                                             |
| default | 10                                                                                                                               |
| flag    | --saa.eval-period-factor                                                                                                         |
| env     | SK_SAA_EVAL_PERIOD_FACTOR                                                                                                        |

### Cleanup Cycle

|         |                                                                   |
| ------- | ----------------------------------------------------------------- |
| name    | cleanup-cycle                                                     |
| usage   | Cycle/ frequency the ScaleAlertAggregator removes expired alerts. |
| type    | Duration                                                          |
| default | 60s                                                               |
| flag    | --saa.cleanup-cycle                                               |
| env     | SK_SAA_CLEANUP_CYCLE                                              |

### Scale Alerts

|         |                                                                                                                                          |
| ------- | ---------------------------------------------------------------------------------------------------------------------------------------- |
| name    | scale-alerts                                                                                                                             |
| usage   | Cycle/ frequency the ScaleAlertAggregator removes expired alerts.                                                                        |
| type    | List of value triplets (alert-name:alert-weight:alert-description). List elements are separated by a ';' and values are separated by '.' |
| default | ""                                                                                                                                       |
| example | --saa.scale-alerts="AlertA:1.0:An upscaling alert;AlertB:-1.5:A downscaling alert"                                                       |
| flag    | --saa.scale-alerts                                                                                                                       |
| env     | SK_SAA_SCALE_ALERTS                                                                                                                      |