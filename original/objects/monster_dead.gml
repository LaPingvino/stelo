// object monster_dead (id 7)
// sprite=sprite_monsterflat solid=False visible=True depth=0 persistent=False parent=none mask=none

#event Create
action_set_alarm(10, 0);

#event Alarm 0
action_kill_object();
