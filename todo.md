# To-Do
# In progress
- [ ] Add web version of game of life
    - [ ] Host in github pages
    - [ ] Maybe add a runner that deploys it whenever there's changes
    - [x] Added basic web support
- [ ] Remove extra stuff from tcell.go
- [ ] Recover UI colors
# Done
- [x] Make shrinking the grid also shrink the background color and redraw the whole screen
- [x] Optimize drawing logic to only check state around live tiles
- [x] Make it possible to click and drag the mouse to draw many tiles
- [X] Gotta make drawing logic and state logic separate. Now we're doing a bit of whatever everywhere and it's causing a bunch of extra draws for nothing when moving mouse over the play area
- [X] Add UI which says play speed, pause/unpause, in the top left corner
- [X] Make it so that map can be drawn even under the UI

