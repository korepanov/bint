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
intnum;
boolfirst_file;
stringcondition;
voidinit(){print("");
num=0;
first_file=True;
root_source="benv/while_program.b";
SET_SOURCE(root_source);
SET_DEST("benv/for_program.b");
};
print("");
voidfinish(){print("");
UNSET_SOURCE();
UNSET_DEST();
};
print("");
boolis_for(stringcommand){print("");
stacks;
stringbuf;
intpos;
buf="for(";
s=ops(command,buf);
s.pop(buf);
if(NOT("end"==buf)){pos=int(buf);
if(NOT(0==pos)){println("for: ERROR");
exit(1);
};
print("");
};
print("");
return(NOT("end"==buf));
};
print("");
voidswitch_files(){print("");
finish();
[print(""),(first_file),goto(#first_end)];
SET_SOURCE("benv/for_program.b");
SET_DEST("benv/for_program2.b");
first_file=False;
goto(#switch_files_e);
#first_end:print("");
SET_SOURCE("benv/for_program2.b");
SET_DEST("benv/for_program.b");
first_file=True;
#switch_files_e:print("");
print("");
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
DEL_DEST("benv/for_program2.b");
DEL_DEST(root_source);
};
print("");
voidmain(){print("");
init();
intcounter;
intinternal_counter;
intcommand_len;
stringsnum;
stringbuf;
stringinc;
intpos;
boolwas_for;
boolwas_internal_for;
was_for=False;
was_internal_for=False;
#next:print("");
switch_command();
if(NOT("end"==command)){if(is_for(command)){was_for=True;
was_internal_for=False;
command_len=len(command);
command=command[4:command_len];
send_command(command);
switch_command();
send_command(command);
switch_command();
snum=str(num);
buf=(("#for"+snum)+":print(\"\")");
send_command(buf);
buf=(("if("+command)+"){print(\"\")");
condition=buf;
send_command(buf);
switch_command();
counter=block_end();
pos=index(command,"{");
if(-1==pos){println("for: ERROR");
exit(1);
};
print("");
pos=(pos-1);
inc=command[0:pos];
#next_internal:print("");
switch_command();
if(COMMAND_COUNTER<counter){if(is_for(command)){was_internal_for=True;
send_command(command);
switch_command();
send_command(command);
switch_command();
send_command(command);
switch_command();
send_command(command);
internal_counter=block_end();
switch_command();
};
print("");
if(was_internal_for){if(COMMAND_COUNTER>internal_counter){was_internal_for=False;
};
print("");
};
print("");
if(("break"==command)AND(NOT(was_internal_for))){command=(("goto(#for"+snum)+"_end)");
};
print("");
if(("continue"==command)AND(NOT(was_internal_for))){send_command(inc);
command=(("goto(#for"+snum)+")");
};
print("");
send_command(command);
goto(#next_internal);
};
print("");
send_command(inc);
send_command(command);
send_command(condition);
buf=(("goto(#for"+snum)+")");
send_command(buf);
buf="}";
send_command(buf);
command=(("#for"+snum)+"_end:print(\"\")");
send_command(command);
num=(num+1);
}else{send_command(command);
print("");
};
print("");
goto(#next);
};
print("");
if(was_for){was_for=False;
switch_files();
COMMAND_COUNTER=0;
goto(#next);
};
print("");
clear_files();
};
print("");
main();
