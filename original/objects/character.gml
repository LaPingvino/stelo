// object character (id 0)
// sprite=sprite_right solid=False visible=True depth=-1 persistent=False parent=none mask=none

#event Create
ammo = 0;
netusxu = 0;

#event Step (normal)
background_hspeed = -1;
if action_if_variable(sprite_index, sprite_left, 0)
background_hspeed = 1;
// Check whether in the air
if action_if_empty(0, 1, 0)  /*relative*/
action_set_gravity(270, 0.5);
else
action_set_gravity(270, 0);
// Limit the vertical speed
if action_if_variable(vspeed, 12, 2)
vspeed = 12;
// Check whether on a ladder
if action_if_object(ladder, 0, 0)  /*relative*/
{
    action_set_gravity(270, 0);
    vspeed = 0;
    action_set_sprite(sprite_climbing, 1);
}
if action_if(sprite_index==sprite_ladder)
// --- code action ---
if (x==xprevious && y==yprevious) {
 x=(round(x/32)*32);
 y=(round(y/32)*32);
}

// --- end code ---
// --- code action ---
if (room_persistent) {
 room_persistent=false;
 sound_stop_all();
 sound_loop(muziko1);
}

if (string_count("sekvosio",keyboard_string)) > 0 {
 keyboard_string="";
 room_goto_next();
}
if (string_count("vivuplilonge",keyboard_string)) > 0 {
 keyboard_string="";
 lives+=1;
}
if (string_count("netusxumin",keyboard_string)) > 0 {
 keyboard_string="";
 netusxu=1;
}
if (string_count("mibedauxras",keyboard_string)) > 0 {
 keyboard_string="";
 netusxu=0;
 sound_stop_all();
 sound_loop(muziko1);
}
if (string_count("esperanto",keyboard_string)) > 0 {
 keyboard_string="";
 sound_stop_all();
 sound_loop(LaEspero);
}
if (string_count("fekalfajro",keyboard_string)) > 0 {
 keyboard_string="";
 with (death) {
  instance_change(block,false);
 }
}
if (string_count("miamasGrandanFraton",keyboard_string)) > 0 {
 keyboard_string="";
 score+=500;
}
if (string_count("pasispasinteco",keyboard_string)) > 0 {
 keyboard_string="";
 highscore_clear();
 score=-1;
}
if (string_count("pliajpoentoj",keyboard_string)) > 0 {
 keyboard_string="";
 room_persistent=true;
 global.ludnivelo=room;
 sound_stop_all();
 sound_loop(bonussono);
 room_goto(bonusnivelo);
}
// --- end code ---

#event Collision with iru_al_bonuso
// --- code action ---
room_persistent=true;
global.ludnivelo=room;
with (other) {instance_destroy();}
sound_stop_all();
sound_loop(bonussono);
room_goto(bonusnivelo);
// --- end code ---

#event Collision with b_monster
if action_if_variable(netusxu, 1, 0)
exit;
if action_if(vspeed > 0 && y < other.y+8)
{
    // Destroy the monster
    action_sound(killmonster, 0);
    action_set_score(10)  /*relative*/;
    with (other) action_move('000010000', 0);
    with (other) action_change_object(monster_dead, 1);
}
else
{
    // Destroy the character
    action_sound(killcharacter, 0);
    action_sleep(1000, 1);
    // --- code action ---
    room_goto(global.ludnivelo);
    // --- end code ---
}

#event Collision with transirejo
if action_if_variable(netusxu, 1, 0)
exit;
with (other) if action_if(used)
exit;
// --- code action ---
xstart=other.x
ystart=other.y
other.used=true;
// --- end code ---
with (other) action_set_sprite(transirejo_a, 1);

#event Collision with bonus
if action_if_variable(netusxu, 1, 0)
exit;
action_sound(finishlevel, 0);
action_set_life(1)  /*relative*/;
with (other) action_kill_object();

#event Collision with death
if action_if_variable(netusxu, 1, 0)
exit;
action_sound(killcharacter, 0);
action_sleep(1000, 1);
action_set_life(-1)  /*relative*/;
action_move_start();

#event Collision with levelexit
if action_if_variable(netusxu, 1, 0)
exit;
action_sound(finishlevel, 0);
action_sleep(1000, 1);
if action_if_next_room()
{
    action_next_room(0);
}

#event Collision with stelo_poentoj
if action_if_variable(netusxu, 1, 0)
exit;
action_sound(getmushroom, 0);
action_set_score(10)  /*relative*/;
with (other) action_kill_object();

#event Collision with flyer
if action_if_variable(netusxu, 1, 0)
exit;
action_sound(killcharacter, 0);
action_sleep(1000, 1);
if action_if(room==bonusnivelo)
{
    // --- code action ---
    room_goto(global.ludnivelo);
    // --- end code ---
}
else
{
    action_set_life(-1)  /*relative*/;
    action_move_start();
}

#event Collision with monster
if action_if_variable(netusxu, 1, 0)
exit;
if action_if(vspeed > 0 && y < other.y+8)
{
    // Destroy the monster
    action_sound(killmonster, 0);
    action_set_score(50)  /*relative*/;
    with (other) action_move('000010000', 0);
    with (other) action_change_object(monster_dead, 1);
}
else
{
    // Destroy the character
    action_sound(killcharacter, 0);
    action_sleep(1000, 1);
    action_set_life(-1)  /*relative*/;
    action_move_start();
}

#event Collision with block
action_move_contact(direction, 12, 0);
action_set_vspeed(0);

#event Keyboard <down>
if action_if_object(ladder, 0, 0)  /*relative*/
{
    if action_if_empty(0, 3, 0)  /*relative*/
    action_move_to(0, 3)  /*relative*/;
}

#event Keyboard <right>
action_set_sprite(sprite_right, 1);
if action_if_empty(4, 0, 0)  /*relative*/
action_move_to(4, 0)  /*relative*/;
else
{
    if action_if_empty(4, -8, 0)  /*relative*/
    {
        action_move_to(4, -8)  /*relative*/;
        action_move_contact(270, 8, 0);
    }
}

#event Keyboard <up>
if action_if_object(ladder, 0, 0)  /*relative*/
{
    if action_if_empty(0, -3, 0)  /*relative*/
    action_move_to(0, -3)  /*relative*/;
}
else
{
    if action_if_collision(0, 1, 0)  /*relative*/
    action_set_vspeed(-10);
}

#event Keyboard <left>
action_set_sprite(sprite_left, 1);
if action_if_empty(-4, 0, 0)  /*relative*/
action_move_to(-4, 0)  /*relative*/;
else
{
    if action_if_empty(-4, -8, 0)  /*relative*/
    {
        action_move_to(-4, -8)  /*relative*/;
        action_move_contact(270, 8, 0);
    }
}

#event Other: game end
action_end_sound(muziko1);

#event Other: outside room
if action_if_variable(y, room_height, 2)
{
    action_sound(killcharacter, 0);
    action_sleep(1000, 1);
    action_move_start();
}

#event KeyPress <space>
if action_if_variable(ammo, 0, 2)
action_create_object(<object 32>, 0, 0)  /*relative*/;
ammo += -1;
