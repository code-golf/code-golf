FROM node:24.2.0-bookworm-slim

RUN apt-get update && apt-get install --no-install-recommends -y \
    ca-certificates curl fontforge fonttools make python3-fontforge unzip woff2 zip

WORKDIR /twemoji-colr

RUN curl -#L https://github.com/matrix-org/twemoji-colr/tarball/397dec7 \
  | tar xz --strip-components 1

RUN npm install

RUN echo [] > extras/ligatures.json                      \
 && unzip -q -j -d svg twe-svg.zip                       \
    '*/1f1*-1f1*.svg'   `# Flags`                        \
    '*/1f0cf.svg'       `# Joker`                        \
    '*/1f308.svg'       `# Rainbow`                      \
    '*/1f30c.svg'       `# Milky Way`                    \
    '*/1f340.svg'       `# Four Leaf Clover`             \
    '*/1f356.svg'       `# Meat on Bone`                 \
    '*/1f35e.svg'       `# Bread`                        \
    '*/1f36f.svg'       `# Honey Pot`                    \
    '*/1f371.svg'       `# Bento Box`                    \
    '*/1f377.svg'       `# Wine Glass`                   \
    '*/1f37a.svg'       `# Beer Mug`                     \
    '*/1f382.svg'       `# Birthday Cake`                \
    '*/1f385.svg'       `# Santa Claus`                  \
    '*/1f388.svg'       `# Balloon`                      \
    '*/1f3ae.svg'       `# Video Game Controller`        \
    '*/1f3af.svg'       `# Direct Hit`                   \
    '*/1f3b3.svg'       `# Bowling`                      \
    '*/1f3b8.svg'       `# Guitar`                       \
    '*/1f3ba.svg'       `# Trumpet`                      \
    '*/1f3c3.svg'       `# Person Running`               \
    '*/1f3c6.svg'       `# Trophy`                       \
    '*/1f3cc.svg'       `# Person Golfing`               \
    '*/1f3de.svg'       `# National Park`                \
    '*/1f3e5.svg'       `# Hospital`                     \
    '*/1f409.svg'       `# Dragon`                       \
    '*/1f40b.svg'       `# Whale`                        \
    '*/1f40d.svg'       `# Sanke`                        \
    '*/1f411.svg'       `# Sheep`                        \
    '*/1f418.svg'       `# Elephant`                     \
    '*/1f41f.svg'       `# Fish`                         \
    '*/1f426.svg'       `# Bird`                         \
    '*/1f42a.svg'       `# Camel`                        \
    '*/1f441.svg'       `# Eye`                          \
    '*/1f444.svg'       `# Mouth`                        \
    '*/1f445.svg'       `# Tongue`                       \
    '*/1f44b.svg'       `# Waving Hand`                  \
    '*/1f44d.svg'       `# Thumbs Up`                    \
    '*/1f47f.svg'       `# Angry Face with Horns`        \
    '*/1f48d.svg'       `# Ring`                         \
    '*/1f48e.svg'       `# Gem Stone`                    \
    '*/1f4af.svg'       `# Hundred Points`               \
    '*/1f4b0.svg'       `# Money Bag`                    \
    '*/1f4bc.svg'       `# Briefcase`                    \
    '*/1f4be.svg'       `# Floppy Disk`                  \
    '*/1f4c4.svg'       `# Page Facing Up`               \
    '*/1f4c5.svg'       `# Calendar`                     \
    '*/1f4d0.svg'       `# Triangular Ruler`             \
    '*/1f4d5.svg'       `# Closed Book`                  \
    '*/1f4da.svg'       `# Books`                        \
    '*/1f4e3.svg'       `# Megaphone`                    \
    '*/1f4f4.svg'       `# Mobile Phone Off`             \
    '*/1f511.svg'       `# Key`                          \
    '*/1f51e.svg'       `# No One Under Eighteen`        \
    '*/1f520.svg'       `# Input Latin Uppercase`        \
    '*/1f523.svg'       `# Input Symbols`                \
    '*/1f526.svg'       `# Flashlight`                   \
    '*/1f549.svg'       `# Om`                           \
    '*/1f596.svg'       `# Vulcan Salute`                \
    '*/1f5a5.svg'       `# Desktop Computer`             \
    '*/1f5d1.svg'       `# Wastebasket`                  \
    '*/1f5dc.svg'       `# Compression`                  \
    '*/1f5f3.svg'       `# Ballot Box with Ballot`       \
    '*/1f600.svg'       `# Grinning Face`                \
    '*/1f602.svg'       `# Face with Tears of Joy`       \
    '*/1f605.svg'       `# Grinning Face with Sweat`     \
    '*/1f606.svg'       `# Grinning Squinting Face`      \
    '*/1f607.svg'       `# Smiling Face with Halo`       \
    '*/1f608.svg'       `# Smiling Face with Horns`      \
    '*/1f609.svg'       `# Winking Face`                 \
    '*/1f60e.svg'       `# Smiling Face with Sunglasses` \
    '*/1f60f.svg'       `# Smirking Face`                \
    '*/1f610.svg'       `# Neutral Face`                 \
    '*/1f613.svg'       `# Downcast Face with Sweat`     \
    '*/1f615.svg'       `# Confused Face`                \
    '*/1f617.svg'       `# Kissing Face`                 \
    '*/1f61b.svg'       `# Face with Tongue`             \
    '*/1f61c.svg'       `# Winking Face with Tongue`     \
    '*/1f61d.svg'       `# Squinting Face with Tongue`   \
    '*/1f621.svg'       `# Pouting Face`                 \
    '*/1f622.svg'       `# Crying Face`                  \
    '*/1f62e.svg'       `# Face with Open Mouth`         \
    '*/1f633.svg'       `# Flushed Face`                 \
    '*/1f634.svg'       `# Sleeping Face`                \
    '*/1f636.svg'       `# Face Without Mouth`           \
    '*/1f641.svg'       `# Slightly Frowning Face`       \
    '*/1f642.svg'       `# Slightly Smiling Face`        \
    '*/1f697.svg'       `# Automobile`                   \
    '*/1f6a2.svg'       `# Ship`                         \
    '*/1f910.svg'       `# Zipper-Mouth Face`            \
    '*/1f947.svg'       `# 1st Place Medal`              \
    '*/1f948.svg'       `# 2nd Place Medal`              \
    '*/1f949.svg'       `# 3rd Place Medal`              \
    '*/1f961.svg'       `# Takeout Box`                  \
    '*/1f963.svg'       `# Bowl with Spoon`              \
    '*/1f967.svg'       `# Pie`                          \
    '*/1f96a.svg'       `# Sandwich`                     \
    '*/1f970.svg'       `# Smiling Face with Hearts`     \
    '*/1f984.svg'       `# Unicorn`                      \
    '*/1f98b.svg'       `# Butterfly`                    \
    '*/1f98e.svg'       `# Lizard`                       \
    '*/1f9a5.svg'       `# Sloth`                        \
    '*/1f9ab.svg'       `# Beaver`                       \
    '*/1f9d9.svg'       `# Mage`                         \
    '*/1f9db.svg'       `# Vampire`                      \
    '*/1f9e0.svg'       `# Brain`                        \
    '*/1f9ea.svg'       `# Test Tube`                    \
    '*/1f9f6.svg'       `# Yarn`                         \
    '*/1fa84.svg'       `# Magic Wand`                   \
    '*/1fa9b.svg'       `# Screwdriver`                  \
    '*/1fa9e.svg'       `# Mirror`                       \
    '*/1faa6.svg'       `# Headstone`                    \
    '*/23f1.svg'        `# Stopwatch`                    \
    '*/2615.svg'        `# Hot Beverage`                 \
    '*/2623.svg'        `# Biohazard`                    \
    '*/2639.svg'        `# Frowning Face`                \
    '*/26f3.svg'        `# Flag in Hole`                 \
    '*/2702.svg'        `# Scissors`                     \
    '*/274c.svg'        `# Ballot X`                     \
    '*/2b50.svg'        `# Star`                         \
 && rm twe-svg.zip                                       \
    svg/1f1ea-1f1fa.svg `# Flag: European Union`         \
    svg/1f1fa-1f1f3.svg `# Flag: United Nations`         \
 && zip -qr twe-svg.zip svg                              \
 && make                                                 \
 && woff2_compress 'build/Twemoji Mozilla.ttf'
