// object monsterfarilo (id 39)
// sprite=sprite_level solid=False visible=False depth=0 persistent=False parent=none mask=none

#event Create
timelineaktiva = 0;

#event Step (normal)
if action_if_number(b_monster, 0, 0)
if action_if(timelineaktiva==0)
{
    action_set_timeline(faruMonsterojn, 0);
    timelineaktiva = 1;
}
