# To-Do
# In progress
- [ ] Make unit tests
- [ ] Pausing causes some inconsistencies in tick timing, especially on slower tick intervals. Setup slow-ish interval and unpause and pause quickly multiple times to reproduce
# Done
- [x] Remove extra stuff from tcell.go
- [x] Recover UI colors
- [x] Clean up code a bit
- [x] Add more keybindings and update them to readme
- [x] Add web version of game of life
    - [x] Host in github pages
    - [x] Maybe add a runner that deploys it whenever there's changes
    - [x] Add a link to github pages to readme
    - [x] Added basic web support
- [x] Make shrinking the grid also shrink the background color and redraw the whole screen
- [x] Optimize drawing logic to only check state around live tiles
- [x] Make it possible to click and drag the mouse to draw many tiles
- [X] Gotta make drawing logic and state logic separate. Now we're doing a bit of whatever everywhere and it's causing a bunch of extra draws for nothing when moving mouse over the play area
- [X] Add UI which says play speed, pause/unpause, in the top left corner
- [X] Make it so that map can be drawn even under the UI

