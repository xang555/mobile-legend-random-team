# Configuration Reference

Configuration defaults live in `configs/config.yaml`. Each field may be overridden with an environment variable using the `RMT_` prefix and replacing dots with underscores (implemented via Viper's key replacer).

## Server
| Key | Env | Description | Default |
| --- | --- | --- | --- |
| `server.host` | `RMT_SERVER_HOST` | Listener host/IP | `0.0.0.0` |
| `server.port` | `RMT_SERVER_PORT` | Listener port | `8080` |
| `server.read_timeout` | `RMT_SERVER_READ_TIMEOUT` | Request read timeout | `5s` |
| `server.write_timeout` | `RMT_SERVER_WRITE_TIMEOUT` | Response write timeout | `5s` |
| `server.idle_timeout` | `RMT_SERVER_IDLE_TIMEOUT` | Keep-alive timeout | `60s` |
| `server.shutdown_timeout` | `RMT_SERVER_SHUTDOWN_TIMEOUT` | Graceful shutdown deadline | `10s` |

## Logging
| Key | Env | Description | Default |
| --- | --- | --- | --- |
| `logging.level` | `RMT_LOGGING_LEVEL` | Log level (`debug`, `info`, `warn`, `error`) | `info` |
| `logging.encoding` | `RMT_LOGGING_ENCODING` | Encoder (`json`, `console`) | `json` |

## Team
| Key | Env | Description | Default |
| --- | --- | --- | --- |
| `team.composition` | `RMT_TEAM_COMPOSITION` | Ordered list of roles to fill. Comma-separated when using env var. | `[Tank, Fighter, Assassin, Mage, Marksman]` |
| `team.allow_duplicates` | `RMT_TEAM_ALLOW_DUPLICATES` | Allow repeated heroes when true | `false` |
| `team.heroes.<role>` | `RMT_TEAM_HEROES_<ROLE>` | Hero pool for each role. Provide comma-separated list via env | Defined per role |

### Overriding arrays via environment variables
Set comma-separated values:
```bash
export RMT_TEAM_COMPOSITION="Tank,Fighter,Support"
export RMT_TEAM_HEROES_SUPPORT="Angela,Estes"
```

## Configuration Validation
The application validates that every role in `team.composition` has at least one hero.
