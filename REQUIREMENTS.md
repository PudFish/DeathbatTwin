# Requirements

## Definitions

**Deathbat:** One of 10,000 unique images from the [Deathbats Club](https://avengedsevenfold.io/).

**Trait:** The characteristics associated with a Deathbat.

**Twin:** Two Deathbats with the most alike traits.

**1/1:** One of One, a truly unique Deathbat that deviates from the standard image and traits and therefore has no twin.

## Functional

**FR01.** The user shall be able to enter a source Deathbat number.

**FR02.** The source Deathbat image shall be displayed.

**FR03.** Deathbat traits shall be displayed along with any Deathbat images.

**FR04.** The source Deathbat traits shall be compared against all Deathbat traits to find their twin.

**FR05.** 1/1 Deathbats shall be determined not to have a twin.

**FR06.** The Deathbat twin shall be displayed.

**FR07.** Current Deathbat owners shall be displayed along with any Deathbat image.

**FR08.** Deathbats with multiple potential twins shall allocate the Deathbat with the closest number to the source Deathbat as the twin.

**FR09.** Deathbats with multiple potential twins equally close shall allocate the lowest number as the twin.

**FR10.** Deathbat numbers shall be displayed along with any Deathbat images.

## Performance

**PR01.** Deathbat traits shall be cached to minimise API calls.

**PR02.** Deathbat images shall be downloaded at 600x600 pixels to minimise network traffic.

**PR03.** Deathbat owners shall be determined at run time to ensure the latest owner is displayed.

**PR04.** Deathbat image URLs shall be cached to minimise API calls.

**PR05.** Deathbat twin pairs shall not be cached to account for changing traits or additional Deathbats created. 
