// object b_monster (id 38)
// sprite=sprite_monsterr solid=False visible=True depth=0 persistent=False parent=none mask=none

#event Create
action_move('000001000', 4);

#event Step (end)
if action_if_variable(hspeed, 0, 2)
action_set_sprite(sprite_monsterr, 1);
else
action_set_sprite(sprite_monsterl, 1);

#event Step (normal)
// Check whether in the air
if action_if_empty(0, 1, 0)  /*relative*/
action_set_gravity(270, 0.5);
else
action_set_gravity(270, 0);

#event Collision with marker
action_move_to(xprevious, yprevious);
action_reverse_xdir();

#event Collision with block
action_move_contact(direction, 12, 0);
action_set_vspeed(0);
action_reverse_xdir();
