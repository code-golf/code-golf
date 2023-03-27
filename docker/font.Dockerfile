FROM node:10.24.1-buster-slim

RUN apt-get update && apt-get install --no-install-recommends -y \
    ca-certificates fontforge fonttools git make python-fontforge unzip woff2 zip

RUN git clone -b v0.6.0 https://github.com/mozilla/twemoji-colr.git

WORKDIR twemoji-colr

RUN npm install

RUN echo [] > extras/ligatures.json                      \
 && unzip -q twe-svg.zip                                 \
    svg/1f1*-*          `# Flags`                        \
    svg/1f0cf.svg       `# Joker`                        \
    svg/1f308.svg       `# Rainbow`                      \
    svg/1f30c.svg       `# Milky Way`                    \
    svg/1f340.svg       `# Four Leaf Clover`             \
    svg/1f356.svg       `# Meat on Bone`                 \
    svg/1f35e.svg       `# Bread`                        \
    svg/1f36f.svg       `# Honey Pot`                    \
    svg/1f377.svg       `# Wine Glass`                   \
    svg/1f37a.svg       `# Beer Mug`                     \
    svg/1f382.svg       `# Birthday Cake`                \
    svg/1f385.svg       `# Santa Claus`                  \
    svg/1f3af.svg       `# Direct Hit`                   \
    svg/1f3b3.svg       `# Bowling`                      \
    svg/1f3b8.svg       `# Guitar`                       \
    svg/1f3ba.svg       `# Trumpet`                      \
    svg/1f3c3.svg       `# Person Running`               \
    svg/1f3c6.svg       `# Trophy`                       \
    svg/1f3e5.svg       `# Hospital`                     \
    svg/1f409.svg       `# Dragon`                       \
    svg/1f40b.svg       `# Whale`                        \
    svg/1f40d.svg       `# Sanke`                        \
    svg/1f418.svg       `# Elephant`                     \
    svg/1f41f.svg       `# Fish`                         \
    svg/1f426.svg       `# Bird`                         \
    svg/1f42a.svg       `# Camel`                        \
    svg/1f441.svg       `# Eye`                          \
    svg/1f445.svg       `# Tongue`                       \
    svg/1f44b.svg       `# Waving Hand`                  \
    svg/1f44d.svg       `# Thumbs Up`                    \
    svg/1f47f.svg       `# Angry Face with Horns`        \
    svg/1f48d.svg       `# Ring`                         \
    svg/1f48e.svg       `# Gem Stone`                    \
    svg/1f4b0.svg       `# Money Bag`                    \
    svg/1f4bc.svg       `# Briefcase`                    \
    svg/1f4be.svg       `# Floppy Disk`                  \
    svg/1f4c4.svg       `# Page Facing Up`               \
    svg/1f4c5.svg       `# Calendar`                     \
    svg/1f4d5.svg       `# Closed Book`                  \
    svg/1f4da.svg       `# Books`                        \
    svg/1f4e3.svg       `# Megaphone`                    \
    svg/1f4f4.svg       `# Mobile Phone Off`             \
    svg/1f511.svg       `# Key`                          \
    svg/1f51e.svg       `# No One Under Eighteen`        \
    svg/1f520.svg       `# Input Latin Uppercase`        \
    svg/1f523.svg       `# Input Symbols`                \
    svg/1f549.svg       `# Om`                           \
    svg/1f596.svg       `# Vulcan Salute`                \
    svg/1f5a5.svg       `# Desktop Computer`             \
    svg/1f5dc.svg       `# Compression`                  \
    svg/1f5f3.svg       `# Ballot Box with Ballot`       \
    svg/1f600.svg       `# Grinning Face`                \
    svg/1f602.svg       `# Face with Tears of Joy`       \
    svg/1f605.svg       `# Grinning Face with Sweat`     \
    svg/1f606.svg       `# Grinning Squinting Face`      \
    svg/1f607.svg       `# Smiling Face with Halo`       \
    svg/1f608.svg       `# Smiling Face with Horns`      \
    svg/1f609.svg       `# Winking Face`                 \
    svg/1f60e.svg       `# Smiling Face with Sunglasses` \
    svg/1f60f.svg       `# Smirking Face`                \
    svg/1f610.svg       `# Neutral Face`                 \
    svg/1f613.svg       `# Downcast Face with Sweat`     \
    svg/1f615.svg       `# Confused Face`                \
    svg/1f617.svg       `# Kissing Face`                 \
    svg/1f61b.svg       `# Face with Tongue`             \
    svg/1f61c.svg       `# Winking Face with Tongue`     \
    svg/1f61d.svg       `# Squinting Face with Tongue`   \
    svg/1f621.svg       `# Pouting Face`                 \
    svg/1f622.svg       `# Crying Face`                  \
    svg/1f62e.svg       `# Face with Open Mouth`         \
    svg/1f633.svg       `# Flushed Face`                 \
    svg/1f634.svg       `# Sleeping Face`                \
    svg/1f636.svg       `# Face Without Mouth`           \
    svg/1f641.svg       `# Slightly Frowning Face`       \
    svg/1f642.svg       `# Slightly Smiling Face`        \
    svg/1f697.svg       `# Automobile`                   \
    svg/1f6a2.svg       `# Ship`                         \
    svg/1f910.svg       `# Zipper-Mouth Face`            \
    svg/1f947.svg       `# 1st Place Medal`              \
    svg/1f948.svg       `# 2nd Place Medal`              \
    svg/1f949.svg       `# 3rd Place Medal`              \
    svg/1f961.svg       `# Takeout Box`                  \
    svg/1f967.svg       `# Pie`                          \
    svg/1f98b.svg       `# Butterfly`                    \
    svg/1f98e.svg       `# Lizard`                       \
    svg/1f9a5.svg       `# Sloth`                        \
    svg/1f9db.svg       `# Vampire`                      \
    svg/1f9e0.svg       `# Brain`                        \
    svg/1f9ea.svg       `# Test Tube`                    \
    svg/1f9f6.svg       `# Yarn`                         \
    svg/1fa9b.svg       `# Screwdriver`                  \
    svg/1fa9e.svg       `# Mirror`                       \
    svg/1faa6.svg       `# Headstone`                    \
    svg/2615.svg        `# Hot Beverage`                 \
    svg/2639.svg        `# Frowning Face`                \
    svg/26f3.svg        `# Flag in Hole`                 \
    svg/2702.svg        `# Scissors`                     \
    svg/2b50.svg        `# Star`                         \
 && rm twe-svg.zip                                       \
    svg/1f1ea-1f1fa.svg `# Flag: European Union`         \
    svg/1f1fa-1f1f3.svg `# Flag: United Nations`         \
 && zip -qr twe-svg.zip svg                              \
 && make                                                 \
 && woff2_compress 'build/Twemoji Mozilla.ttf'
