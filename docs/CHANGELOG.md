| Release | Date        | Comments                                                                                                           |
|---------|-------------|--------------------------------------------------------------------------------------------------------------------|
| 1.8.0   | 2026.07.25  | Upgraded to latest customError and helperFunctions<br>Version number now SemVer-aligned<br>BuildDeps & GO updates  |
| 1.61.01 | 2025.10.15  | GO version bump, cosmetic README.md updates                                                                        |
| 1.61.00 | 2025.08.25  | GO version bump                                                                                                    |
| 1.60.00 | 2025.06.06  | Updated to GO 1.24.4 to incorporate the crypto/x509 bugfixes                                                       |
| 1.59.00 | 2025.01.13  | CN will default to the certificate name if omitted                                                                 |
| 1.58.00 | 2024.11.08  | Fixed missing config path issue                                                                                    |
| 1.57.00 | 2024.07.09  | Reverted code to 1.55.00, added better cerr handling                                                               |
| 1.56.00 | 2024.06.25  | Reverted to durations in years, topping at 2 for CAs                                                               | 
| 1.55.02 | 2024.06.07  | Missed one certification expiration in years, updated to GO 1.22.4                                                 |
| 1.55.01 | 2024.06.02  | Certificate expiration time is now computed in days, not years                                                     |
| 1.52.00 | 2024.05.11  | Fixed error handling part two...                                                                                   |
| 1.51.00 | 2024.05.10  | Fixed error handling where functions natively used customError                                                     |
| 1.50.00 | 2024.05.03  | Most helper functions + custom errors moved to outside packages                                                    |
| 1.25.01 | 2024.03.28  | Minor : CERTDEFAULTS and CADEFAULTS can be in lowercase                                                            |         
| 1.25.00 | 2024.03.28  | Moved all configs from $HOME/.config/ to $HOME/.config/JFG/                                                        |
| 1.24.xx | 2024.03.28  | UNPUBLISHED : minor output fixes, dep updates                                                                      |
| 1.23.00 | 2023.12.07  | GO version bump to 1.21.5, package updates                                                                         |
| 1.22.00 | 2023.11.02  | fixed broken cert rm (duplicated path in pathname                                                                  |
| 1.21.00 | 2023.11.02  | defaultEnv.json does not need to be specified anymore in cm env                                                    |
| 1.20.06 | 2023.10.31  | Error message in cert verify less confusing                                                                        |
| 1.20.05 | 2023.10.16  | New version numbering scheme                                                                                       |
| 1.205   | 2023.10.15  | Fixed wrong path for server's private keys                                                                         |
| 1.200   | 2023.10.13  | Go version bump to 1.21.3, simplified path handling in environment                                                 |
| 1.100   | 2023.10.04  | `cm cert verify` now displays comments when asked to                                                               |
| 1.010-1 | 2023.10.03  | Serial number now properly incremented                                                                             |
| 1.001   | 2023.10.03  | More verbosity, doc update                                                                                         |
| 1.000   | 2023.09.30  | Prod-ready version.                                                                                                |
| 0.100   | 2023.08.13  | Initial version.                                                                                                   |




