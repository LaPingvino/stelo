// object klarigilo (id 41)
// sprite=none solid=False visible=True depth=0 persistent=False parent=none mask=none

#event Other: room start
// --- code action ---
screen_redraw();
keyboard_wait();
room_goto_next();
// --- end code ---

#event Other: game start
action_sound(komenco, 1);
// --- code action ---
{
 depth_malnov=depth;
 depth=-2;
 draw_background(stelo_komenco,0,0);
 screen_refresh();
 keyboard_wait();
 depth=depth_malnov;
 global.startita=false;
}
// --- end code ---
action_end_sound(komenco);
action_sound(muziko1, 1);
