// object life_controller (id 34)
// sprite=none solid=False visible=True depth=-1 persistent=True parent=none mask=none

#event Other: no more lives
action_another_room(finlisto, 0);

#event Other: room start
if action_if(global.startita)
exit;
action_set_life(3);
action_set_caption(0, 'score: ', 0, 'lives: ', 0, 'health: ');
action_draw_font("@Arial Unicode MS",10,0,0,0,0,0);
global.startita = true;

#event Draw
action_color(11855591, 0);
action_draw_rectangle(view_left+5, view_top+5, view_left+100, view_top+25);
action_draw_life_images(character.x, character.y-10, sprite_life);
action_draw_score(view_left+10, view_top+5, 'Poentoj: ');
