intCOMMAND_COUNTER;
stringcommand;
COMMAND_COUNTER=0;
stackreverse(stacks){print("");
stringbuf;
stackres;
s.pop(buf);
#_reverse_s:print("");
[goto(#_reverse_e),("end"==buf),print("")];
res.push(buf);
s.pop(buf);
goto(#_reverse_s);
#_reverse_e:print("");
returnres;
};
print("");
stackindexes(strings,stringsub_s){print("");
stackres;
inti;
intpointer;
ints_len;
ints_len_old;
intsub_len;
s_len=len(s);
sub_len=len(sub_s);
s_len_old=s_len;
i=index(s,sub_s);
pointer=i;
#_indexes_s:print("");
[goto(#_indexes_e),(-1==i),print("")];
i=(i+(s_len_old-s_len));
res.push(i);
pointer=(pointer+sub_len);
s=s[pointer:s_len];
s_len=len(s);
i=index(s,sub_s);
pointer=i;
goto(#_indexes_s);
#_indexes_e:print("");
res=reverse(res);
returnres;
};
print("");
voidSET_COMMAND_COUNTER(intcounter){print("");
inti;
i=0;
stringcommand;
RESET_SOURCE();
COMMAND_COUNTER=counter;
#_set_start:print("");
[print(""),(i<counter),goto(#_set_end)];
next_command(command);
i=(i+1);
goto(#_set_start);
#_set_end:print("");
print("");
};
print("");
stringget_command(intcounter){print("");
inti;
stringbuf;
stringcommand;
i=0;
RESET_SOURCE();
#_get_command_s:print("");
[print(""),(i<counter),goto(#_get_command_e)];
next_command(command);
i=(i+1);
goto(#_get_command_s);
#_get_command_e:print("");
SET_COMMAND_COUNTER(COMMAND_COUNTER);
returncommand;
};
print("");
voidswitch_command(){print("");
COMMAND_COUNTER=(COMMAND_COUNTER+1);
next_command(command);
};
print("");
intstack_len(stacks){print("");
intres;
stringbuf;
res=0;
s.pop(buf);
#_stack_len_s:print("");
[goto(#_stack_len_e),("end"==buf),print("")];
res=(res+1);
s.pop(buf);
goto(#_stack_len_s);
#_stack_len_e:print("");
returnres;
};
print("");
stackops(stringcommand,stringop){print("");
stackquotes;
stackthese_quotes;
stackop_nums;
stringbuf;
stackres;
intnum1;
intnum2;
intop_num;
boolwas_quote;
boolto_add;
was_quote=False;
to_add=True;
op_nums=indexes(command,op);
op_nums.pop(buf);
#_op_nums_s:print("");
[goto(#_op_nums_e),("end"==buf),print("")];
op_num=int(buf);
quotes=reg_find("\"(\\.|[^\"])*\"",command);
#_quotes_s:print("");
quotes.pop(these_quotes);
these_quotes.pop(buf);
[goto(#_quotes_e),("end"==buf),print("")];
#_these_quotes_s:print("");
[goto(#_these_quotes_e),("end"==buf),print("")];
num1=int(buf);
these_quotes.pop(buf);
num2=int(buf);
these_quotes.pop(buf);
was_quote=True;
goto(#_these_quotes_s);
#_these_quotes_e:print("");
[print(""),((op_num>num1)AND(op_num<num2)),goto(#_is_op_end)];
to_add=False;
goto(#_push_op_end);
#_is_op_end:print("");
goto(#_quotes_s);
#_quotes_e:print("");
[goto(#_push_op_end),((was_quote)AND(NOT(to_add))),print("")];
res.push(op_num);
was_quote=False;
to_add=True;
#_push_op_end:print("");
print("");
op_nums.pop(buf);
was_quote=False;
to_add=True;
goto(#_op_nums_s);
#_op_nums_e:print("");
res=reverse(res);
returnres;
};
print("");
intblock_end(){print("");
stringop1;
stringop2;
stringcode;
stringcommand_buf;
into_sum;
intc_sum;
intcommand_len;
stackobraces;
stackcbraces;
stringbuf;
stringspos;
intcounter;
intbuf_counter;
intpos;
counter=COMMAND_COUNTER;
command_len=len(command);
command=command[1:command_len];
code=command;
next_command(command);
counter=(counter+1);
code=(code+command[0]);
o_sum=1;
c_sum=0;
#_block_s:print("");
[goto(#_block_e),(o_sum==c_sum),print("")];
obraces=ops(code,"{");
cbraces=ops(code,"}");
o_sum=stack_len(obraces);
c_sum=stack_len(cbraces);
command_len=len(command);
command_buf=command[1:command_len];
code=(code+command_buf);
next_command(command);
counter=(counter+1);
command_buf=command[0];
code=(code+command_buf);
goto(#_block_s);
#_block_e:print("");
buf_counter=(COMMAND_COUNTER-1);
SET_COMMAND_COUNTER(buf_counter);
next_command(command);
COMMAND_COUNTER=(COMMAND_COUNTER+1);
return(counter-1);
};
print("");
voidprintln(strings){print("");
print(s);
print("\n");
};
print("");
boolin_stack(stacks,stringel){print("");
stringbuf;
boolres;
res=False;
s.pop(buf);
#_in_stack_s:print("");
[goto(#_in_stack_e),("end"==buf),print("")];
[print(""),(el==buf),goto(#_no)];
res=True;
goto(#_in_stack_e);
#_no:print("");
s.pop(buf);
goto(#_in_stack_s);
#_in_stack_e:print("");
returnres;
};
print("");
intfunc_end(stringcommand,intfunc_begin){print("");
stackobraces;
stackcbraces;
stringobrace;
stringcbrace;
stringsymbol;
into_sum;
intc_sum;
intpos;
stringspos;
intcommand_len;
command_len=len(command);
obrace="(";
cbrace=")";
o_sum=1;
c_sum=0;
pos=(func_begin+1);
obraces=ops(command,obrace);
cbraces=ops(command,cbrace);
#_braces_s:print("");
[print(""),(pos<command_len),goto(#_braces_e)];
spos=str(pos);
[print(""),in_stack(obraces,spos),goto(#_o_plus_end)];
o_sum=(o_sum+1);
#_o_plus_end:print("");
[print(""),in_stack(cbraces,spos),goto(#_c_plus_end)];
c_sum=(c_sum+1);
#_c_plus_end:print("");
[goto(#_braces_e),(o_sum==c_sum),print("")];
pos=(pos+1);
goto(#_braces_s);
#_braces_e:print("");
returnpos;
};
print("");
stringroot_source;
stringcommand;
boolfirst_file;
intnum;
intexit_num;
intbr_closed;
intbr_opened;
voidinit(){print("");
br_closed=0;
br_opened=0;
root_source="benv/for_program.b";
SET_SOURCE(root_source);
SET_DEST("benv/if_program.b");
};
print("");
voidfinish(){print("");
UNSET_SOURCE();
UNSET_DEST();
};
print("");
boolis_var_def(stringcommand){print("");
stacks;
stringbuf;
s=reg_find("^(?:(int|float|bool|stack|string)[^\(]*)",command);
s.pop(buf);
return(NOT("end"==buf));
};
print("");
stringType(stringcommand){print("");
stacks;
stringbuf;
s=reg_find("(?:(^int))",command);
s.pop(buf);
[goto(#int_end),("end"==buf),print("")];
return"int";
#int_end:print("");
s=reg_find("(?:(^float))",command);
s.pop(buf);
[goto(#float_end),("end"==buf),print("")];
return"float";
#float_end:print("");
s=reg_find("(?:(^bool))",command);
s.pop(buf);
[goto(#bool_end),("end"==buf),print("")];
return"bool";
#bool_end:print("");
s=reg_find("(?:(^stack))",command);
s.pop(buf);
[goto(#stack_end),("end"==buf),print("")];
return"stack";
#stack_end:print("");
s=reg_find("(?:(^string))",command);
s.pop(buf);
[goto(#string_end),("end"==buf),print("")];
return"string";
#string_end:print("");
s=reg_find("(?:(^void))",command);
s.pop(buf);
[goto(#void_end),("end"==buf),print("")];
return"void";
#void_end:print("");
print(command);
print("\n");
print("Type: ERROR\n");
exit(1);
};
print("");
voidcheck_br(stringcommand){print("");
stringsymbol;
stacks;
stringbuf;
symbol="{";
s=ops(command,symbol);
s.pop(buf);
[goto(#br_opened_end),("end"==buf),print("")];
br_opened=(br_opened+1);
#br_opened_end:print("");
symbol="}";
s=ops(command,symbol);
s.pop(buf);
[goto(#br_closed_end),("end"==buf),print("")];
br_closed=(br_closed+1);
#br_closed_end:print("");
print("");
};
print("");
voidreset_br(){print("");
br_opened=0;
br_closed=0;
};
print("");
boolis_if(stringcommand){print("");
stringop;
stacks;
stringbuf;
intibuf;
intstart_pos;
op="if(";
s=ops(command,op);
s.pop(buf);
[goto(#end_true),("end"==buf),print("")];
ibuf=int(buf);
[goto(#start_true),(0==ibuf),print("")];
start_pos=(ibuf-4);
[goto(#end_true),(start_pos<0),print("")];
buf=command[start_pos:ibuf];
[print(""),("else"==buf),goto(#end_true)];
#start_true:print("");
returnTrue;
#end_true:print("");
returnFalse;
print("");
};
print("");
stringget_cond(stringcommand){print("");
stringop;
stacks;
stringbuf;
intstart_pos;
intend_pos;
op="if(";
s=ops(command,op);
s.pop(buf);
start_pos=int(buf);
start_pos=(start_pos+2);
end_pos=func_end(command,start_pos);
end_pos=(end_pos+1);
buf=command[start_pos:end_pos];
returnbuf;
};
print("");
stringif_type(stringcommand){print("");
intcommand_len;
stringprefix;
command_len=len(command);
[print(""),(1==command_len),goto(#not_clear)];
return"clear";
#not_clear:print("");
prefix=command[1:7];
[print(""),("elseif"==prefix),goto(#not_elseif)];
return"elseif";
#not_elseif:print("");
prefix=command[1:5];
[print(""),("else"==prefix),goto(#if_type_error)];
return"else";
#if_type_error:print("");
return"error";
};
print("");
voidswitch_files(){print("");
finish();
[print(""),(first_file),goto(#first_end)];
SET_SOURCE("benv/if_program.b");
SET_DEST("benv/if_program2.b");
first_file=False;
goto(#switch_files_e);
#first_end:print("");
SET_SOURCE("benv/if_program2.b");
SET_DEST("benv/if_program.b");
first_file=True;
#switch_files_e:print("");
print("");
};
print("");
voidreplace_if(stringcond,intstop_pos){print("");
stringbuf;
stringsnum;
stackargs_to_undefine;
stackargs_to_undefine_old;
stringarg_type;
intcommand_len;
stringarg_name;
boolwas_br;
snum=str(num);
buf=(((("[print(\"\"), "+cond)+", goto(#_cond")+snum)+"_end)]");
send_command(buf);
switch_command();
check_br(command);
[print(""),(br_opened>(br_closed+1)),goto(#was_not_br)];
was_br=True;
#was_not_br:print("");
print("");
#replace_clear_if_s:print("");
[goto(#replace_clear_if_e),("end"==command),print("")];
[print(""),(stop_pos==COMMAND_COUNTER),goto(#add_replace_clear_if_mark)];
args_to_undefine.pop(arg_name);
#un:print("");
[goto(#un_end),("end"==arg_name),print("")];
buf=(("UNDEFINE("+arg_name)+")");
send_command(buf);
args_to_undefine.pop(arg_name);
goto(#un);
#un_end:print("");
buf=(("#_cond"+snum)+"_end:print(\"\")");
send_command(buf);
switch_command();
check_br(command);
[print(""),(br_opened>(br_closed+1)),goto(#was_not_br2)];
was_br=True;
#was_not_br2:print("");
print("");
#add_replace_clear_if_mark:print("");
[print(""),(len(command)>6),goto(#replace_if_ret_end)];
print("");
[print(""),(command[0:6]=="return"),goto(#replace_if_ret_end)];
print("");
[print(""),(COMMAND_COUNTER<stop_pos),goto(#replace_if_ret_end)];
args_to_undefine_old=args_to_undefine;
args_to_undefine.pop(arg_name);
#un6:print("");
[goto(#un_end6),("end"==arg_name),print("")];
buf=(("UNDEFINE("+arg_name)+")");
send_command(buf);
args_to_undefine.pop(arg_name);
goto(#un6);
#un_end6:print("");
print("");
args_to_undefine=args_to_undefine_old;
#replace_if_ret_end:print("");
[print(""),((is_var_def(command))AND(br_closed==br_opened)),goto(#pop_e)];
arg_type=Type(command);
inttype_len;
type_len=len(arg_type);
command_len=len(command);
arg_name=command[type_len:command_len];
args_to_undefine.push(arg_name);
#pop_e:print("");
send_command(command);
switch_command();
check_br(command);
[print(""),(br_opened>(br_closed+1)),goto(#was_not_br3)];
was_br=True;
#was_not_br3:print("");
print("");
goto(#replace_clear_if_s);
#replace_clear_if_e:print("");
reset_br();
COMMAND_COUNTER=0;
switch_files();
};
print("");
voidreplace_elseif(stringcond,intstop_pos){print("");
print("");
stringbuf;
intibuf;
stringsnum;
stringt;
stringbcommand;
intthis_stop_pos;
stringsexit_num;
intcommand_len;
stringarg_name;
stringarg_type;
stackargs_to_undefine;
stackargs_to_undefine_old;
sexit_num=str(exit_num);
snum=str(num);
buf=(((("[print(\"\"), "+cond)+", goto(#_cond")+snum)+"_end)]");
send_command(buf);
switch_command();
check_br(command);
#replace_elseif_s:print("");
[goto(#replace_elseif_e),("end"==command),print("")];
[print(""),(stop_pos==COMMAND_COUNTER),goto(#add_replace_elseif_mark)];
args_to_undefine.pop(arg_name);
#un2:print("");
[goto(#un_end2),("end"==arg_name),print("")];
buf=(("UNDEFINE("+arg_name)+")");
send_command(buf);
args_to_undefine.pop(arg_name);
goto(#un2);
#un_end2:print("");
reset_br();
sexit_num=str(exit_num);
buf=(("goto(#_cond_exit"+sexit_num)+")");
send_command(buf);
buf=(("#_cond"+snum)+"_end:print(\"\")");
send_command(buf);
num=(num+1);
t=if_type(command);
[print(""),("elseif"==t),goto(#find_block_e)];
this_stop_pos=block_end();
#find_block_e:print("");
bcommand=command;
switch_command();
check_br(command);
[goto(#replace_elseif_e),(NOT("elseif"==t)),print("")];
[print(""),("elseif"==t),goto(#elif_end)];
cond=get_cond(bcommand);
snum=str(num);
buf=(((("[print(\"\"), "+cond)+", goto(#_cond")+snum)+"_end)]");
send_command(buf);
stop_pos=this_stop_pos;
goto(#replace_elseif_s);
#elif_end:print("");
print("");
#add_replace_elseif_mark:print("");
[print(""),(len(command)>6),goto(#replace_else_if_ret_end)];
print("");
[print(""),(command[0:6]=="return"),goto(#replace_else_if_ret_end)];
print("");
[print(""),(COMMAND_COUNTER<stop_pos),goto(#replace_else_if_ret_end)];
args_to_undefine_old=args_to_undefine;
args_to_undefine.pop(arg_name);
#un7:print("");
[goto(#un_end7),("end"==arg_name),print("")];
buf=(("UNDEFINE("+arg_name)+")");
send_command(buf);
args_to_undefine.pop(arg_name);
goto(#un7);
#un_end7:print("");
print("");
args_to_undefine=args_to_undefine_old;
#replace_else_if_ret_end:print("");
print("");
[print(""),((is_var_def(command))AND(br_closed==br_opened)),goto(#pop_e2)];
arg_type=Type(command);
inttype_len;
type_len=len(arg_type);
command_len=len(command);
arg_name=command[type_len:command_len];
args_to_undefine.push(arg_name);
#pop_e2:print("");
send_command(command);
switch_command();
check_br(command);
goto(#replace_elseif_s);
#replace_elseif_e:print("");
t=if_type(bcommand);
[print(""),("else"==t),goto(#restore_end)];
ibuf=(COMMAND_COUNTER-2);
SET_COMMAND_COUNTER(ibuf);
switch_command();
stop_pos=block_end();
switch_command();
check_br(command);
#restore_end:print("");
[print(""),("else"==t),goto(#final_cond_end)];
#final_cond_s:print("");
[goto(#final_cond_end),(stop_pos==COMMAND_COUNTER),print("")];
print("");
[print(""),((is_var_def(command))AND(br_closed==br_opened)),goto(#pop_e3)];
arg_type=Type(command);
inttype_len;
type_len=len(arg_type);
command_len=len(command);
arg_name=command[type_len:command_len];
args_to_undefine.push(arg_name);
#pop_e3:print("");
send_command(command);
switch_command();
check_br(command);
goto(#final_cond_s);
#final_cond_end:print("");
[print(""),("else"==t),goto(#else_end)];
switch_command();
check_br(command);
#else_end:print("");
args_to_undefine.pop(arg_name);
#un3:print("");
[goto(#un_end3),("end"==arg_name),print("")];
buf=(("UNDEFINE("+arg_name)+")");
send_command(buf);
args_to_undefine.pop(arg_name);
goto(#un3);
#un_end3:print("");
reset_br();
sexit_num=str(exit_num);
buf=(("#_cond_exit"+sexit_num)+":print(\"\")");
exit_num=(exit_num+1);
send_command(buf);
#ts:print("");
[goto(#te),("end"==command),print("")];
send_command(command);
switch_command();
goto(#ts);
#te:print("");
COMMAND_COUNTER=0;
reset_br();
switch_files();
};
print("");
voidreplace_else(stringcond,intstop_pos){print("");
stringbuf;
stringsnum;
stringsexit_num;
intibuf;
intpos;
intcommand_len;
stringarg_name;
stringarg_type;
stackargs_to_undefine;
pos=-1;
snum=str(num);
buf=(((("[print(\"\"), "+cond)+", goto(#_cond")+snum)+"_end)]");
send_command(buf);
switch_command();
check_br(command);
#replace_else_s:print("");
[goto(#replace_else_e),("end"==command),print("")];
[print(""),(stop_pos==COMMAND_COUNTER),goto(#add_replace_else_mark)];
sexit_num=str(exit_num);
args_to_undefine.pop(arg_name);
#un4:print("");
[goto(#un_end4),("end"==arg_name),print("")];
buf=(("UNDEFINE("+arg_name)+")");
send_command(buf);
args_to_undefine.pop(arg_name);
goto(#un4);
#un_end4:print("");
reset_br();
buf=(("goto(#_cond_exit"+sexit_num)+")");
send_command(buf);
buf=(("#_cond"+snum)+"_end:print(\"\")");
send_command(buf);
switch_command();
ibuf=(COMMAND_COUNTER-2);
SET_COMMAND_COUNTER(ibuf);
switch_command();
pos=block_end();
switch_command();
check_br(command);
#add_replace_else_mark:print("");
[print(""),(pos==COMMAND_COUNTER),goto(#figure_brace_end)];
args_to_undefine.pop(arg_name);
#un5:print("");
[goto(#un_end5),("end"==arg_name),print("")];
buf=(("UNDEFINE("+arg_name)+")");
send_command(buf);
args_to_undefine.pop(arg_name);
goto(#un5);
#un_end5:print("");
sexit_num=str(exit_num);
buf=(("#_cond_exit"+sexit_num)+":print(\"\")");
exit_num=(exit_num+1);
send_command(buf);
switch_command();
#ets:print("");
[goto(#ete),("end"==command),print("")];
send_command(command);
switch_command();
goto(#ets);
#ete:print("");
goto(#replace_else_e);
#figure_brace_end:print("");
[print(""),((is_var_def(command))AND(br_closed==br_opened)),goto(#pop_e4)];
arg_type=Type(command);
inttype_len;
type_len=len(arg_type);
command_len=len(command);
arg_name=command[type_len:command_len];
args_to_undefine.push(arg_name);
#pop_e4:print("");
send_command(command);
switch_command();
check_br(command);
goto(#replace_else_s);
#replace_else_e:print("");
COMMAND_COUNTER=0;
reset_br();
switch_files();
};
print("");
voidclear_files(){print("");
[goto(#first_file),(first_file),print("")];
switch_files();
switch_command();
#clear_files_s:print("");
[goto(#clear_files_e),("end"==command),print("")];
send_command(command);
switch_command();
goto(#clear_files_s);
#first_file:print("");
print("");
#clear_files_e:print("");
finish();
DEL_DEST("benv/if_program2.b");
DEL_DEST(root_source);
};
print("");
voidmain(){print("");
stringbuf;
stringcond;
intcounter;
intbuf_counter;
stringsnum;
stringt;
boolis_stop;
first_file=True;
init();
num=0;
exit_num=0;
COMMAND_COUNTER=0;
#again_s:print("");
is_stop=True;
switch_command(command);
#main_s:print("");
[goto(#main_e),("end"==command),print("")];
[print(""),(is_if(command)),goto(#next)];
is_stop=False;
cond=get_cond(command);
counter=block_end();
buf=get_command(counter);
t=if_type(buf);
[print(""),("error"==t),goto(#error_end)];
println("ERROR in the if type\n");
goto(#total_e);
#error_end:print("");
[print(""),("clear"==t),goto(#if_end)];
replace_if(cond,counter);
num=(num+1);
goto(#main_e);
#if_end:print("");
[print(""),("elseif"==t),goto(#elseif_end)];
replace_elseif(cond,counter);
goto(#main_e);
#elseif_end:print("");
replace_else(cond,counter);
num=(num+1);
goto(#main_e);
#next:print("");
send_command(command);
switch_command();
goto(#main_s);
#main_e:print("");
[print(""),(is_stop),goto(#again_s)];
#total_e:print("");
clear_files();
};
print("");
main();
