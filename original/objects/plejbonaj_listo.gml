// object plejbonaj_listo (id 37)
// sprite=none solid=False visible=True depth=0 persistent=False parent=none mask=none

#event Create
with (life_controller) action_kill_object();
// --- code action ---
sound_stop_all();
sound_loop(finavenko);
// --- end code ---
action_draw_font("Arial Rondigita MT Grasa ",20,65280,0,0,0,0);

#event Step (normal)
// --- code action ---
if (pluiru) {
 keyboard_wait();
 sound_stop_all();
 game_restart();
}
// --- end code ---

#event Other: room start
// --- code action ---
if (score > highscore_value(10)) {
 keyboard_string='';
 while (ord(keyboard_lastchar)<>13){
  draw_background(plejbonaj,0,0);
  draw_text(70,70,'Vi tiom bone ludis ke vi aperos en#la listo de la plej bonaj ludantoj!##Via nomo: '+keyboard_string+'_');
  screen_refresh();
  io_handle();
 }
 highscore_add(keyboard_string,score);
}
else
{
 draw_background(plejbonaj,0,0);
 draw_text(70,70,'Ho ve... Vi ne suficxe bone ludis por#aperi en la listo de la plej bonaj#ludantoj...');
 screen_refresh();
 keyboard_wait();
}
// --- end code ---

#event Draw
// --- code action ---
draw_highscore(x+5,y+5,x+480,y+320);
pluiru=1;
// --- end code ---
