$filestdlib/core.b$;
$trace1$;
intCOMMAND_COUNTER;
$trace2$;
stringcommand;
$trace3$;
COMMAND_COUNTER=0;
$trace4$;
stackreverse(stacks){stringbuf;
$trace5$;
stackres;
$trace6$;
s.pop(buf);
$trace7$;
#_reverse_s:[goto(#_reverse_e),("end"==buf),print("")];
$trace8$;
res.push(buf);
$trace9$;
s.pop(buf);
$trace10$;
goto(#_reverse_s);
$trace11$;
#_reverse_e:returnres;
$trace12$;
};
$trace13$;
stackindexes(strings,stringsub_s){stackres;
$trace14$;
inti;
$trace15$;
intpointer;
$trace16$;
ints_len;
$trace17$;
ints_len_old;
$trace18$;
intsub_len;
$trace19$;
s_len=len(s);
$trace20$;
sub_len=len(sub_s);
$trace21$;
s_len_old=s_len;
$trace22$;
i=index(s,sub_s);
$trace23$;
pointer=i;
$trace24$;
#_indexes_s:[goto(#_indexes_e),(-1==i),print("")];
$trace25$;
i=(i+(s_len_old-s_len));
$trace26$;
res.push(i);
$trace27$;
pointer=(pointer+sub_len);
$trace28$;
s=s[pointer:s_len];
$trace29$;
s_len=len(s);
$trace30$;
i=index(s,sub_s);
$trace31$;
pointer=i;
$trace32$;
goto(#_indexes_s);
$trace33$;
#_indexes_e:res=reverse(res);
$trace34$;
returnres;
$trace35$;
};
$trace36$;
voidSET_COMMAND_COUNTER(intcounter){inti;
$trace37$;
i=0;
$trace38$;
stringcommand;
$trace39$;
RESET_SOURCE();
$trace40$;
COMMAND_COUNTER=counter;
$trace41$;
#_set_start:[print(""),(i<counter),goto(#_set_end)];
$trace42$;
next_command(command);
$trace43$;
i=(i+1);
$trace44$;
goto(#_set_start);
$trace45$;
#_set_end:print("");
$trace46$;
};
$trace47$;
stringget_command(intcounter){inti;
$trace48$;
stringbuf;
$trace49$;
stringcommand;
$trace50$;
i=0;
$trace51$;
RESET_SOURCE();
$trace52$;
#_get_command_s:[print(""),(i<counter),goto(#_get_command_e)];
$trace53$;
next_command(command);
$trace54$;
i=(i+1);
$trace55$;
goto(#_get_command_s);
$trace56$;
#_get_command_e:SET_COMMAND_COUNTER(COMMAND_COUNTER);
$trace57$;
returncommand;
$trace58$;
};
$trace59$;
voidswitch_command(){COMMAND_COUNTER=(COMMAND_COUNTER+1);
$trace60$;
next_command(command);
$trace61$;
};
$trace62$;
intstack_len(stacks){intres;
$trace63$;
stringbuf;
$trace64$;
res=0;
$trace65$;
s.pop(buf);
$trace66$;
#_stack_len_s:[goto(#_stack_len_e),("end"==buf),print("")];
$trace67$;
res=(res+1);
$trace68$;
s.pop(buf);
$trace69$;
goto(#_stack_len_s);
$trace70$;
#_stack_len_e:returnres;
$trace71$;
};
$trace72$;
stackops(stringcommand,stringop){stackquotes;
$trace73$;
stackthese_quotes;
$trace74$;
stackop_nums;
$trace75$;
stringbuf;
$trace76$;
stackres;
$trace77$;
intnum1;
$trace78$;
intnum2;
$trace79$;
intop_num;
$trace80$;
boolwas_quote;
$trace81$;
boolto_add;
$trace82$;
was_quote=False;
$trace83$;
to_add=True;
$trace84$;
op_nums=indexes(command,op);
$trace85$;
op_nums.pop(buf);
$trace86$;
#_op_nums_s:[goto(#_op_nums_e),("end"==buf),print("")];
$trace87$;
op_num=int(buf);
$trace88$;
quotes=reg_find("\"(\\.|[^\"])*\"",command);
$trace89$;
#_quotes_s:quotes.pop(these_quotes);
$trace90$;
these_quotes.pop(buf);
$trace91$;
[goto(#_quotes_e),("end"==buf),print("")];
$trace92$;
#_these_quotes_s:[goto(#_these_quotes_e),("end"==buf),print("")];
$trace93$;
num1=int(buf);
$trace94$;
these_quotes.pop(buf);
$trace95$;
num2=int(buf);
$trace96$;
these_quotes.pop(buf);
$trace97$;
was_quote=True;
$trace98$;
goto(#_these_quotes_s);
$trace99$;
#_these_quotes_e:[print(""),((op_num>num1)AND(op_num<num2)),goto(#_is_op_end)];
$trace100$;
to_add=False;
$trace101$;
goto(#_push_op_end);
$trace102$;
#_is_op_end:goto(#_quotes_s);
$trace103$;
#_quotes_e:[goto(#_push_op_end),((was_quote)AND(NOT(to_add))),print("")];
$trace104$;
res.push(op_num);
$trace105$;
was_quote=False;
$trace106$;
to_add=True;
$trace107$;
#_push_op_end:print("");
$trace108$;
op_nums.pop(buf);
$trace109$;
was_quote=False;
$trace110$;
to_add=True;
$trace111$;
goto(#_op_nums_s);
$trace112$;
#_op_nums_e:res=reverse(res);
$trace113$;
returnres;
$trace114$;
};
$trace115$;
intblock_end(){stringop1;
$trace116$;
stringop2;
$trace117$;
stringcode;
$trace118$;
stringcommand_buf;
$trace119$;
into_sum;
$trace120$;
intc_sum;
$trace121$;
intcommand_len;
$trace122$;
stackobraces;
$trace123$;
stackcbraces;
$trace124$;
stringbuf;
$trace125$;
stringspos;
$trace126$;
intcounter;
$trace127$;
intbuf_counter;
$trace128$;
intpos;
$trace129$;
counter=COMMAND_COUNTER;
$trace130$;
command_len=len(command);
$trace131$;
command=command[1:command_len];
$trace132$;
code=command;
$trace133$;
next_command(command);
$trace134$;
counter=(counter+1);
$trace135$;
code=(code+command[0]);
$trace136$;
o_sum=1;
$trace137$;
c_sum=0;
$trace138$;
#_block_s:[goto(#_block_e),(o_sum==c_sum),print("")];
$trace139$;
obraces=ops(code,"{");
$trace140$;
cbraces=ops(code,"}");
$trace141$;
o_sum=stack_len(obraces);
$trace142$;
c_sum=stack_len(cbraces);
$trace143$;
command_len=len(command);
$trace144$;
command_buf=command[1:command_len];
$trace145$;
code=(code+command_buf);
$trace146$;
next_command(command);
$trace147$;
counter=(counter+1);
$trace148$;
command_buf=command[0];
$trace149$;
code=(code+command_buf);
$trace150$;
goto(#_block_s);
$trace151$;
#_block_e:buf_counter=(COMMAND_COUNTER-1);
$trace152$;
SET_COMMAND_COUNTER(buf_counter);
$trace153$;
next_command(command);
$trace154$;
COMMAND_COUNTER=(COMMAND_COUNTER+1);
$trace155$;
return(counter-1);
$trace156$;
};
$trace157$;
voidprintln(strings){print(s);
$trace158$;
print("\n");
$trace159$;
};
$trace160$;
boolin_stack(stacks,stringel){stringbuf;
$trace161$;
boolres;
$trace162$;
res=False;
$trace163$;
s.pop(buf);
$trace164$;
#_in_stack_s:[goto(#_in_stack_e),("end"==buf),print("")];
$trace165$;
[print(""),(el==buf),goto(#_no)];
$trace166$;
res=True;
$trace167$;
goto(#_in_stack_e);
$trace168$;
#_no:s.pop(buf);
$trace169$;
goto(#_in_stack_s);
$trace170$;
#_in_stack_e:returnres;
$trace171$;
};
$trace172$;
intfunc_end(stringcommand,intfunc_begin){stackobraces;
$trace173$;
stackcbraces;
$trace174$;
stringobrace;
$trace175$;
stringcbrace;
$trace176$;
stringsymbol;
$trace177$;
into_sum;
$trace178$;
intc_sum;
$trace179$;
intpos;
$trace180$;
stringspos;
$trace181$;
intcommand_len;
$trace182$;
command_len=len(command);
$trace183$;
obrace="(";
$trace184$;
cbrace=")";
$trace185$;
o_sum=1;
$trace186$;
c_sum=0;
$trace187$;
pos=(func_begin+1);
$trace188$;
obraces=ops(command,obrace);
$trace189$;
cbraces=ops(command,cbrace);
$trace190$;
#_braces_s:[print(""),(pos<command_len),goto(#_braces_e)];
$trace191$;
spos=str(pos);
$trace192$;
[print(""),in_stack(obraces,spos),goto(#_o_plus_end)];
$trace193$;
o_sum=(o_sum+1);
$trace194$;
#_o_plus_end:[print(""),in_stack(cbraces,spos),goto(#_c_plus_end)];
$trace195$;
c_sum=(c_sum+1);
$trace196$;
#_c_plus_end:[goto(#_braces_e),(o_sum==c_sum),print("")];
$trace197$;
pos=(pos+1);
$trace198$;
goto(#_braces_s);
$trace199$;
#_braces_e:returnpos;
$trace200$;
};
$filebenv/if.b$;
$trace1$;
stringroot_source;
$trace2$;
stringcommand;
$trace3$;
boolfirst_file;
$trace4$;
intnum;
$trace5$;
intexit_num;
$trace6$;
intbr_closed;
$trace7$;
intbr_opened;
$trace8$;
voidinit(){br_closed=0;
$trace9$;
br_opened=0;
$trace10$;
root_source="benv/for_program.b";
$trace11$;
SET_SOURCE(root_source);
$trace12$;
SET_DEST("benv/if_program.b");
$trace13$;
};
$trace14$;
voidfinish(){UNSET_SOURCE();
$trace15$;
UNSET_DEST();
$trace16$;
};
$trace17$;
boolis_var_def(stringcommand){stacks;
$trace18$;
stringbuf;
$trace19$;
s=reg_find("^(?:(int|float|bool|stack|string)[^\(]*)",command);
$trace20$;
s.pop(buf);
$trace21$;
return(NOT("end"==buf));
$trace22$;
};
$trace23$;
stringType(stringcommand){stacks;
$trace24$;
stringbuf;
$trace25$;
s=reg_find("(?:(^int))",command);
$trace26$;
s.pop(buf);
$trace27$;
[goto(#int_end),("end"==buf),print("")];
$trace28$;
return"int";
$trace29$;
#int_end:s=reg_find("(?:(^float))",command);
$trace30$;
s.pop(buf);
$trace31$;
[goto(#float_end),("end"==buf),print("")];
$trace32$;
return"float";
$trace33$;
#float_end:s=reg_find("(?:(^bool))",command);
$trace34$;
s.pop(buf);
$trace35$;
[goto(#bool_end),("end"==buf),print("")];
$trace36$;
return"bool";
$trace37$;
#bool_end:s=reg_find("(?:(^stack))",command);
$trace38$;
s.pop(buf);
$trace39$;
[goto(#stack_end),("end"==buf),print("")];
$trace40$;
return"stack";
$trace41$;
#stack_end:s=reg_find("(?:(^string))",command);
$trace42$;
s.pop(buf);
$trace43$;
[goto(#string_end),("end"==buf),print("")];
$trace44$;
return"string";
$trace45$;
#string_end:s=reg_find("(?:(^void))",command);
$trace46$;
s.pop(buf);
$trace47$;
[goto(#void_end),("end"==buf),print("")];
$trace48$;
return"void";
$trace49$;
#void_end:print(command);
$trace50$;
print("\n");
$trace51$;
print("Type: ERROR\n");
$trace52$;
exit(1);
$trace53$;
};
$trace54$;
voidcheck_br(stringcommand){stringsymbol;
$trace55$;
stacks;
$trace56$;
stringbuf;
$trace57$;
symbol="{";
$trace58$;
s=ops(command,symbol);
$trace59$;
s.pop(buf);
$trace60$;
[goto(#br_opened_end),("end"==buf),print("")];
$trace61$;
br_opened=(br_opened+1);
$trace62$;
#br_opened_end:symbol="}";
$trace63$;
s=ops(command,symbol);
$trace64$;
s.pop(buf);
$trace65$;
[goto(#br_closed_end),("end"==buf),print("")];
$trace66$;
br_closed=(br_closed+1);
$trace67$;
#br_closed_end:print("");
$trace68$;
};
$trace69$;
voidreset_br(){br_opened=0;
$trace70$;
br_closed=0;
$trace71$;
};
$trace72$;
boolis_if(stringcommand){stringop;
$trace73$;
stacks;
$trace74$;
stringbuf;
$trace75$;
intibuf;
$trace76$;
intstart_pos;
$trace77$;
op="if(";
$trace78$;
s=ops(command,op);
$trace79$;
s.pop(buf);
$trace80$;
[goto(#end_true),("end"==buf),print("")];
$trace81$;
ibuf=int(buf);
$trace82$;
[goto(#start_true),(0==ibuf),print("")];
$trace83$;
start_pos=(ibuf-4);
$trace84$;
[goto(#end_true),(start_pos<0),print("")];
$trace85$;
buf=command[start_pos:ibuf];
$trace86$;
[print(""),("else"==buf),goto(#end_true)];
$trace87$;
#start_true:returnTrue;
$trace88$;
#end_true:returnFalse;
$trace89$;
print("");
$trace90$;
};
$trace91$;
stringget_cond(stringcommand){stringop;
$trace92$;
stacks;
$trace93$;
stringbuf;
$trace94$;
intstart_pos;
$trace95$;
intend_pos;
$trace96$;
op="if(";
$trace97$;
s=ops(command,op);
$trace98$;
s.pop(buf);
$trace99$;
start_pos=int(buf);
$trace100$;
start_pos=(start_pos+2);
$trace101$;
end_pos=func_end(command,start_pos);
$trace102$;
end_pos=(end_pos+1);
$trace103$;
buf=command[start_pos:end_pos];
$trace104$;
returnbuf;
$trace105$;
};
$trace106$;
stringif_type(stringcommand){intcommand_len;
$trace107$;
stringprefix;
$trace108$;
command_len=len(command);
$trace109$;
[print(""),(1==command_len),goto(#not_clear)];
$trace110$;
return"clear";
$trace111$;
#not_clear:prefix=command[1:7];
$trace112$;
[print(""),("elseif"==prefix),goto(#not_elseif)];
$trace113$;
return"elseif";
$trace114$;
#not_elseif:prefix=command[1:5];
$trace115$;
[print(""),("else"==prefix),goto(#if_type_error)];
$trace116$;
return"else";
$trace117$;
#if_type_error:return"error";
$trace118$;
};
$trace119$;
voidswitch_files(){finish();
$trace120$;
[print(""),(first_file),goto(#first_end)];
$trace121$;
SET_SOURCE("benv/if_program.b");
$trace122$;
SET_DEST("benv/if_program2.b");
$trace123$;
first_file=False;
$trace124$;
goto(#switch_files_e);
$trace125$;
#first_end:SET_SOURCE("benv/if_program2.b");
$trace126$;
SET_DEST("benv/if_program.b");
$trace127$;
first_file=True;
$trace128$;
#switch_files_e:print("");
$trace129$;
};
$trace130$;
voidreplace_if(stringcond,intstop_pos){stringbuf;
$trace131$;
stringsnum;
$trace132$;
stackargs_to_undefine;
$trace133$;
stringarg_type;
$trace134$;
intcommand_len;
$trace135$;
stringarg_name;
$trace136$;
snum=str(num);
$trace137$;
buf=(((("[print(\"\"), "+cond)+", goto(#_cond")+snum)+"_end)]");
$trace138$;
send_command(buf);
$trace139$;
switch_command();
$trace140$;
#replace_clear_if_s:[goto(#replace_clear_if_e),("end"==command),print("")];
$trace141$;
[print(""),(stop_pos==COMMAND_COUNTER),goto(#add_replace_clear_if_mark)];
$trace142$;
args_to_undefine.pop(arg_name);
$trace143$;
#un:[goto(#un_end),("end"==arg_name),print("")];
$trace144$;
buf=(("UNDEFINE("+arg_name)+")");
$trace145$;
send_command(buf);
$trace146$;
args_to_undefine.pop(arg_name);
$trace147$;
goto(#un);
$trace148$;
#un_end:buf=(("#_cond"+snum)+"_end:print(\"\")");
$trace149$;
send_command(buf);
$trace150$;
#add_replace_clear_if_mark:[print(""),(len(command)>6),goto(#replace_if_ret_end)];
$trace151$;
print("");
$trace152$;
[print(""),(command[0:6]=="return"),goto(#replace_if_ret_end)];
$trace153$;
args_to_undefine.pop(arg_name);
$trace154$;
#un6:[goto(#un_end6),("end"==arg_name),print("")];
$trace155$;
buf=(("UNDEFINE("+arg_name)+")");
$trace156$;
send_command(buf);
$trace157$;
args_to_undefine.pop(arg_name);
$trace158$;
goto(#un6);
$trace159$;
#un_end6:print("");
$trace160$;
#replace_if_ret_end:switch_command();
$trace161$;
check_br(command);
$trace162$;
[print(""),((is_var_def(command))AND(br_closed==br_opened)),goto(#pop_e)];
$trace163$;
arg_type=Type(command);
$trace164$;
inttype_len;
$trace165$;
type_len=len(arg_type);
$trace166$;
command_len=len(command);
$trace167$;
arg_name=command[type_len:command_len];
$trace168$;
args_to_undefine.push(arg_name);
$trace169$;
#pop_e:send_command(command);
$trace170$;
switch_command();
$trace171$;
goto(#replace_clear_if_s);
$trace172$;
#replace_clear_if_e:reset_br();
$trace173$;
COMMAND_COUNTER=0;
$trace174$;
switch_files();
$trace175$;
};
$trace176$;
voidreplace_elseif(stringcond,intstop_pos){stringbuf;
$trace177$;
intibuf;
$trace178$;
stringsnum;
$trace179$;
stringt;
$trace180$;
stringbcommand;
$trace181$;
intthis_stop_pos;
$trace182$;
stringsexit_num;
$trace183$;
intcommand_len;
$trace184$;
stringarg_name;
$trace185$;
stringarg_type;
$trace186$;
stackargs_to_undefine;
$trace187$;
sexit_num=str(exit_num);
$trace188$;
snum=str(num);
$trace189$;
buf=(((("[print(\"\"), "+cond)+", goto(#_cond")+snum)+"_end)]");
$trace190$;
send_command(buf);
$trace191$;
switch_command();
$trace192$;
check_br(command);
$trace193$;
#replace_elseif_s:[goto(#replace_elseif_e),("end"==command),print("")];
$trace194$;
[print(""),(stop_pos==COMMAND_COUNTER),goto(#add_replace_elseif_mark)];
$trace195$;
args_to_undefine.pop(arg_name);
$trace196$;
#un2:[goto(#un_end2),("end"==arg_name),print("")];
$trace197$;
buf=(("UNDEFINE("+arg_name)+")");
$trace198$;
send_command(buf);
$trace199$;
args_to_undefine.pop(arg_name);
$trace200$;
goto(#un2);
$trace201$;
#un_end2:reset_br();
$trace202$;
sexit_num=str(exit_num);
$trace203$;
buf=(("goto(#_cond_exit"+sexit_num)+")");
$trace204$;
send_command(buf);
$trace205$;
buf=(("#_cond"+snum)+"_end:print(\"\")");
$trace206$;
send_command(buf);
$trace207$;
num=(num+1);
$trace208$;
t=if_type(command);
$trace209$;
[print(""),("elseif"==t),goto(#find_block_e)];
$trace210$;
this_stop_pos=block_end();
$trace211$;
#find_block_e:bcommand=command;
$trace212$;
switch_command();
$trace213$;
check_br(command);
$trace214$;
[goto(#replace_elseif_e),(NOT("elseif"==t)),print("")];
$trace215$;
[print(""),("elseif"==t),goto(#elif_end)];
$trace216$;
cond=get_cond(bcommand);
$trace217$;
snum=str(num);
$trace218$;
buf=(((("[print(\"\"), "+cond)+", goto(#_cond")+snum)+"_end)]");
$trace219$;
send_command(buf);
$trace220$;
stop_pos=this_stop_pos;
$trace221$;
goto(#replace_elseif_s);
$trace222$;
#elif_end:print("");
$trace223$;
#add_replace_elseif_mark:print("");
$trace224$;
[print(""),((is_var_def(command))AND(br_closed==br_opened)),goto(#pop_e2)];
$trace225$;
arg_type=Type(command);
$trace226$;
inttype_len;
$trace227$;
type_len=len(arg_type);
$trace228$;
command_len=len(command);
$trace229$;
arg_name=command[type_len:command_len];
$trace230$;
args_to_undefine.push(arg_name);
$trace231$;
#pop_e2:send_command(command);
$trace232$;
switch_command();
$trace233$;
check_br(command);
$trace234$;
goto(#replace_elseif_s);
$trace235$;
#replace_elseif_e:t=if_type(bcommand);
$trace236$;
[print(""),("else"==t),goto(#restore_end)];
$trace237$;
ibuf=(COMMAND_COUNTER-2);
$trace238$;
SET_COMMAND_COUNTER(ibuf);
$trace239$;
switch_command();
$trace240$;
stop_pos=block_end();
$trace241$;
switch_command();
$trace242$;
check_br(command);
$trace243$;
#restore_end:[print(""),("else"==t),goto(#final_cond_end)];
$trace244$;
#final_cond_s:[goto(#final_cond_end),(stop_pos==COMMAND_COUNTER),print("")];
$trace245$;
print("");
$trace246$;
[print(""),((is_var_def(command))AND(br_closed==br_opened)),goto(#pop_e3)];
$trace247$;
arg_type=Type(command);
$trace248$;
inttype_len;
$trace249$;
type_len=len(arg_type);
$trace250$;
command_len=len(command);
$trace251$;
arg_name=command[type_len:command_len];
$trace252$;
args_to_undefine.push(arg_name);
$trace253$;
#pop_e3:send_command(command);
$trace254$;
switch_command();
$trace255$;
check_br(command);
$trace256$;
goto(#final_cond_s);
$trace257$;
#final_cond_end:[print(""),("else"==t),goto(#else_end)];
$trace258$;
switch_command();
$trace259$;
check_br(command);
$trace260$;
#else_end:args_to_undefine.pop(arg_name);
$trace261$;
#un3:[goto(#un_end3),("end"==arg_name),print("")];
$trace262$;
buf=(("UNDEFINE("+arg_name)+")");
$trace263$;
send_command(buf);
$trace264$;
args_to_undefine.pop(arg_name);
$trace265$;
goto(#un3);
$trace266$;
#un_end3:reset_br();
$trace267$;
sexit_num=str(exit_num);
$trace268$;
buf=(("#_cond_exit"+sexit_num)+":print(\"\")");
$trace269$;
exit_num=(exit_num+1);
$trace270$;
send_command(buf);
$trace271$;
#ts:[goto(#te),("end"==command),print("")];
$trace272$;
send_command(command);
$trace273$;
switch_command();
$trace274$;
goto(#ts);
$trace275$;
#te:COMMAND_COUNTER=0;
$trace276$;
reset_br();
$trace277$;
switch_files();
$trace278$;
};
$trace279$;
voidreplace_else(stringcond,intstop_pos){stringbuf;
$trace280$;
stringsnum;
$trace281$;
stringsexit_num;
$trace282$;
intibuf;
$trace283$;
intpos;
$trace284$;
intcommand_len;
$trace285$;
stringarg_name;
$trace286$;
stringarg_type;
$trace287$;
stackargs_to_undefine;
$trace288$;
pos=-1;
$trace289$;
snum=str(num);
$trace290$;
buf=(((("[print(\"\"), "+cond)+", goto(#_cond")+snum)+"_end)]");
$trace291$;
send_command(buf);
$trace292$;
switch_command();
$trace293$;
check_br(command);
$trace294$;
#replace_else_s:[goto(#replace_else_e),("end"==command),print("")];
$trace295$;
[print(""),(stop_pos==COMMAND_COUNTER),goto(#add_replace_else_mark)];
$trace296$;
sexit_num=str(exit_num);
$trace297$;
args_to_undefine.pop(arg_name);
$trace298$;
#un4:[goto(#un_end4),("end"==arg_name),print("")];
$trace299$;
buf=(("UNDEFINE("+arg_name)+")");
$trace300$;
send_command(buf);
$trace301$;
args_to_undefine.pop(arg_name);
$trace302$;
goto(#un4);
$trace303$;
#un_end4:reset_br();
$trace304$;
buf=(("goto(#_cond_exit"+sexit_num)+")");
$trace305$;
send_command(buf);
$trace306$;
buf=(("#_cond"+snum)+"_end:print(\"\")");
$trace307$;
send_command(buf);
$trace308$;
switch_command();
$trace309$;
ibuf=(COMMAND_COUNTER-2);
$trace310$;
SET_COMMAND_COUNTER(ibuf);
$trace311$;
switch_command();
$trace312$;
pos=block_end();
$trace313$;
switch_command();
$trace314$;
check_br(command);
$trace315$;
#add_replace_else_mark:[print(""),(pos==COMMAND_COUNTER),goto(#figure_brace_end)];
$trace316$;
args_to_undefine.pop(arg_name);
$trace317$;
#un5:[goto(#un_end5),("end"==arg_name),print("")];
$trace318$;
buf=(("UNDEFINE("+arg_name)+")");
$trace319$;
send_command(buf);
$trace320$;
args_to_undefine.pop(arg_name);
$trace321$;
goto(#un5);
$trace322$;
#un_end5:sexit_num=str(exit_num);
$trace323$;
buf=(("#_cond_exit"+sexit_num)+":print(\"\")");
$trace324$;
exit_num=(exit_num+1);
$trace325$;
send_command(buf);
$trace326$;
switch_command();
$trace327$;
#ets:[goto(#ete),("end"==command),print("")];
$trace328$;
send_command(command);
$trace329$;
switch_command();
$trace330$;
goto(#ets);
$trace331$;
#ete:goto(#replace_else_e);
$trace332$;
#figure_brace_end:[print(""),((is_var_def(command))AND(br_closed==br_opened)),goto(#pop_e4)];
$trace333$;
arg_type=Type(command);
$trace334$;
inttype_len;
$trace335$;
type_len=len(arg_type);
$trace336$;
command_len=len(command);
$trace337$;
arg_name=command[type_len:command_len];
$trace338$;
args_to_undefine.push(arg_name);
$trace339$;
#pop_e4:send_command(command);
$trace340$;
switch_command();
$trace341$;
check_br(command);
$trace342$;
goto(#replace_else_s);
$trace343$;
#replace_else_e:COMMAND_COUNTER=0;
$trace344$;
reset_br();
$trace345$;
switch_files();
$trace346$;
};
$trace347$;
voidclear_files(){[goto(#first_file),(first_file),print("")];
$trace348$;
switch_files();
$trace349$;
switch_command();
$trace350$;
#clear_files_s:[goto(#clear_files_e),("end"==command),print("")];
$trace351$;
send_command(command);
$trace352$;
switch_command();
$trace353$;
goto(#clear_files_s);
$trace354$;
#first_file:print("");
$trace355$;
#clear_files_e:finish();
$trace356$;
DEL_DEST("benv/if_program2.b");
$trace357$;
DEL_DEST(root_source);
$trace358$;
};
$trace359$;
voidmain(){stringbuf;
$trace360$;
stringcond;
$trace361$;
intcounter;
$trace362$;
intbuf_counter;
$trace363$;
stringsnum;
$trace364$;
stringt;
$trace365$;
boolis_stop;
$trace366$;
first_file=True;
$trace367$;
init();
$trace368$;
num=0;
$trace369$;
exit_num=0;
$trace370$;
COMMAND_COUNTER=0;
$trace371$;
#again_s:is_stop=True;
$trace372$;
switch_command(command);
$trace373$;
#main_s:[goto(#main_e),("end"==command),print("")];
$trace374$;
[print(""),(is_if(command)),goto(#next)];
$trace375$;
is_stop=False;
$trace376$;
cond=get_cond(command);
$trace377$;
counter=block_end();
$trace378$;
buf=get_command(counter);
$trace379$;
t=if_type(buf);
$trace380$;
[print(""),("error"==t),goto(#error_end)];
$trace381$;
println("ERROR in the if type\n");
$trace382$;
goto(#total_e);
$trace383$;
#error_end:[print(""),("clear"==t),goto(#if_end)];
$trace384$;
replace_if(cond,counter);
$trace385$;
num=(num+1);
$trace386$;
goto(#main_e);
$trace387$;
#if_end:[print(""),("elseif"==t),goto(#elseif_end)];
$trace388$;
replace_elseif(cond,counter);
$trace389$;
goto(#main_e);
$trace390$;
#elseif_end:replace_else(cond,counter);
$trace391$;
num=(num+1);
$trace392$;
goto(#main_e);
$trace393$;
#next:send_command(command);
$trace394$;
switch_command();
$trace395$;
goto(#main_s);
$trace396$;
#main_e:[print(""),(is_stop),goto(#again_s)];
$trace397$;
#total_e:clear_files();
$trace398$;
};
$trace399$;
main();
