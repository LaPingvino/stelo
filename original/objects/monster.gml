// object monster (id 5)
// sprite=sprite_monsterr solid=False visible=True depth=0 persistent=False parent=none mask=none

#event Create
action_move('000001000', 4);

#event Step (end)
if action_if_variable(hspeed, 0, 2)
action_set_sprite(sprite_monsterr, 1);
else
action_set_sprite(sprite_monsterl, 1);

#event Collision with marker
action_move_to(xprevious, yprevious);
action_reverse_xdir();

#event Collision with block
action_reverse_xdir();
