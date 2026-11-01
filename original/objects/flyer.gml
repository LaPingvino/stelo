// object flyer (id 8)
// sprite=sprite_flyingr solid=False visible=True depth=0 persistent=False parent=none mask=none

#event Create
action_move('000001000', 3);

#event Step (end)
if action_if_variable(hspeed, 0, 2)
action_set_sprite(sprite_flyingr, 1);
else
action_set_sprite(sprite_flyingl, 1);

#event Collision with marker
action_reverse_xdir();

#event Collision with block
action_reverse_xdir();
