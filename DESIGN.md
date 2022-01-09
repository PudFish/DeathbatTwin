# Design

## User Interface
- Top Centre
  - Title
  - Deathbat number input
- Middle Left
  - Source Deathbat image
  - Source Deathbat number
  - Source Deathbat traits
  - Source Deathbat owner
  - Source Deathbat OpenSea.io hyperlink
- Middle Right
  - Twin Deathbat image
  - Twin Deathbat number
  - Twin Deathbat traits
  - Twin Deathbat owner
  - Twin Deathbat OpenSea.io hyperlink

## APIs
**Traits:** *GET* https://avengedsevenfold.io/deathbats/token/{token_id}

**Image:** *GET* https://avengedsevenfold.io/deathbats/media/{token_id}.jpg

**Owner:** *GET* https://api.opensea.io/api/v1/asset/0x1D3aDa5856B14D9dF178EA5Cab137d436dC55F1D/{token_id}/

**OpenSea Link:** https://opensea.io/assets/0x1d3ada5856b14d9df178ea5cab137d436dc55f1d/{token_id}

## Backend
**Success path:** 
get source Deathbat number --> check source Deathbat number --> retrieve Deathbat traits from cache --> 
compare source Deathbat traits to collection traits --> determine potential twins --> select final twin 
--> get source and twin images --> get source and twin numbers --> get source and twin traits --> 
get source and twin owners --> get source and twin OpenSea.io hyperlink --> display source and twins

**Errors:**

Invalid source Deathbat number

Cache fail

No twin due 1/1

No twin due internal error

API miss

## Trait Twin Priority
Highest to lowest weight

1. 1/1 (Brooks Wackerman, Johnny Christ, M. Shadows, Synyser Gates, Zacky Vengeance)
1. Mask
1. Facial Hair 
1. Eyes, Mouth, Nose
1. Head
1. Skin
1. Background
