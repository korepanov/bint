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
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
$l0=len(s);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
s_len=$l0;
UNDEFINE($l0);
$trace20$;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
$l0=len(sub_s);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
sub_len=$l0;
UNDEFINE($l0);
$trace21$;
s_len_old=s_len;
$trace22$;
int$ind0;
int$ind0;
int$ind0;
int$ind0;
int$ind0;
int$ind0;
int$ind0;
int$ind0;
int$ind0;
int$ind0;
int$ind0;
int$ind0;
int$ind0;
int$ind0;
int$ind0;
int$ind0;
int$ind0;
int$ind0;
int$ind0;
$ind0=index(s,sub_s);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
i=$ind0;
UNDEFINE($ind0);
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
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
$l0=len(s);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
s_len=$l0;
UNDEFINE($l0);
$trace30$;
int$ind0;
int$ind0;
int$ind0;
int$ind0;
int$ind0;
int$ind0;
int$ind0;
int$ind0;
int$ind0;
int$ind0;
int$ind0;
int$ind0;
int$ind0;
int$ind0;
int$ind0;
int$ind0;
int$ind0;
int$ind0;
int$ind0;
$ind0=index(s,sub_s);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
i=$ind0;
UNDEFINE($ind0);
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
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
$I0=int(buf);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
op_num=$I0;
UNDEFINE($I0);
$trace88$;
stack$stack_var0;
stack$stack_var0;
stack$stack_var0;
stack$stack_var0;
stack$stack_var0;
stack$stack_var0;
stack$stack_var0;
stack$stack_var0;
stack$stack_var0;
stack$stack_var0;
stack$stack_var0;
stack$stack_var0;
stack$stack_var0;
stack$stack_var0;
stack$stack_var0;
stack$stack_var0;
stack$stack_var0;
stack$stack_var0;
stack$stack_var0;
$stack_var0=reg_find("\"(\\.|[^\"])*\"",command);
$stack_var0=$stack_var0;
UNDEFINE($stack_var0);
$stack_var0=$stack_var0;
UNDEFINE($stack_var0);
$stack_var0=$stack_var0;
UNDEFINE($stack_var0);
$stack_var0=$stack_var0;
UNDEFINE($stack_var0);
$stack_var0=$stack_var0;
UNDEFINE($stack_var0);
$stack_var0=$stack_var0;
UNDEFINE($stack_var0);
$stack_var0=$stack_var0;
UNDEFINE($stack_var0);
$stack_var0=$stack_var0;
UNDEFINE($stack_var0);
$stack_var0=$stack_var0;
UNDEFINE($stack_var0);
$stack_var0=$stack_var0;
UNDEFINE($stack_var0);
$stack_var0=$stack_var0;
UNDEFINE($stack_var0);
$stack_var0=$stack_var0;
UNDEFINE($stack_var0);
$stack_var0=$stack_var0;
UNDEFINE($stack_var0);
$stack_var0=$stack_var0;
UNDEFINE($stack_var0);
$stack_var0=$stack_var0;
UNDEFINE($stack_var0);
$stack_var0=$stack_var0;
UNDEFINE($stack_var0);
$stack_var0=$stack_var0;
UNDEFINE($stack_var0);
$stack_var0=$stack_var0;
UNDEFINE($stack_var0);
quotes=$stack_var0;
UNDEFINE($stack_var0);
$trace89$;
#_quotes_s:quotes.pop(these_quotes);
$trace90$;
these_quotes.pop(buf);
$trace91$;
[goto(#_quotes_e),("end"==buf),print("")];
$trace92$;
#_these_quotes_s:[goto(#_these_quotes_e),("end"==buf),print("")];
$trace93$;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
$I0=int(buf);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
num1=$I0;
UNDEFINE($I0);
$trace94$;
these_quotes.pop(buf);
$trace95$;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
$I0=int(buf);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
num2=$I0;
UNDEFINE($I0);
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
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
$l0=len(command);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
command_len=$l0;
UNDEFINE($l0);
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
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
$l0=len(command);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
command_len=$l0;
UNDEFINE($l0);
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
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
$l0=len(command);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
command_len=$l0;
UNDEFINE($l0);
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
string$s0;
string$s0;
string$s0;
string$s0;
string$s0;
string$s0;
string$s0;
string$s0;
string$s0;
string$s0;
string$s0;
string$s0;
string$s0;
string$s0;
string$s0;
string$s0;
string$s0;
string$s0;
string$s0;
$s0=str(pos);
$s0=$s0;
UNDEFINE($s0);
$s0=$s0;
UNDEFINE($s0);
$s0=$s0;
UNDEFINE($s0);
$s0=$s0;
UNDEFINE($s0);
$s0=$s0;
UNDEFINE($s0);
$s0=$s0;
UNDEFINE($s0);
$s0=$s0;
UNDEFINE($s0);
$s0=$s0;
UNDEFINE($s0);
$s0=$s0;
UNDEFINE($s0);
$s0=$s0;
UNDEFINE($s0);
$s0=$s0;
UNDEFINE($s0);
$s0=$s0;
UNDEFINE($s0);
$s0=$s0;
UNDEFINE($s0);
$s0=$s0;
UNDEFINE($s0);
$s0=$s0;
UNDEFINE($s0);
$s0=$s0;
UNDEFINE($s0);
$s0=$s0;
UNDEFINE($s0);
$s0=$s0;
UNDEFINE($s0);
spos=$s0;
UNDEFINE($s0);
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
$filebenv/for.b$;
$trace1$;
stringroot_source;
$trace2$;
intnum;
$trace3$;
boolfirst_file;
$trace4$;
voidinit(){num=0;
$trace5$;
first_file=True;
$trace6$;
root_source="benv/while_program.b";
$trace7$;
SET_SOURCE(root_source);
$trace8$;
SET_DEST("benv/for_program.b");
$trace9$;
};
$trace10$;
voidfinish(){UNSET_SOURCE();
$trace11$;
UNSET_DEST();
$trace12$;
};
$trace13$;
boolis_for(stringcommand){stacks;
$trace14$;
stringbuf;
$trace15$;
intpos;
$trace16$;
buf="for(";
$trace17$;
s=ops(command,buf);
$trace18$;
s.pop(buf);
$trace19$;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
int$I0;
$I0=int(buf);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
$I0=$I0;
UNDEFINE($I0);
if(NOT("end"==buf)){pos=$I0;
UNDEFINE($I0);
$trace20$;
if(NOT(0==pos)){println("for: ERROR");
$trace21$;
exit(1);
$trace22$;
};
$trace23$;
};
$trace24$;
return(NOT("end"==buf));
$trace25$;
};
$trace26$;
voidswitch_files(){finish();
$trace27$;
[print(""),(first_file),goto(#first_end)];
$trace28$;
SET_SOURCE("benv/for_program.b");
$trace29$;
SET_DEST("benv/for_program2.b");
$trace30$;
first_file=False;
$trace31$;
goto(#switch_files_e);
$trace32$;
#first_end:SET_SOURCE("benv/for_program2.b");
$trace33$;
SET_DEST("benv/for_program.b");
$trace34$;
first_file=True;
$trace35$;
#switch_files_e:print("");
$trace36$;
};
$trace37$;
voidclear_files(){[goto(#first_file),(first_file),print("")];
$trace38$;
switch_files();
$trace39$;
switch_command();
$trace40$;
#clear_files_s:[goto(#clear_files_e),("end"==command),print("")];
$trace41$;
send_command(command);
$trace42$;
switch_command();
$trace43$;
goto(#clear_files_s);
$trace44$;
#first_file:print("");
$trace45$;
#clear_files_e:finish();
$trace46$;
DEL_DEST("benv/for_program2.b");
$trace47$;
DEL_DEST(root_source);
$trace48$;
};
$trace49$;
voidmain(){init();
$trace50$;
intcounter;
$trace51$;
intinternal_counter;
$trace52$;
intcommand_len;
$trace53$;
stringsnum;
$trace54$;
stringbuf;
$trace55$;
stringinc;
$trace56$;
intpos;
$trace57$;
boolwas_for;
$trace58$;
boolwas_internal_for;
$trace59$;
was_for=False;
$trace60$;
was_internal_for=False;
$trace61$;
#next:switch_command();
$trace62$;
if(NOT("end"==command)){if(is_for(command)){was_for=True;
$trace63$;
was_internal_for=False;
$trace64$;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
int$l0;
$l0=len(command);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
$l0=$l0;
UNDEFINE($l0);
command_len=$l0;
UNDEFINE($l0);
$trace65$;
command=command[4:command_len];
$trace66$;
send_command(command);
$trace67$;
switch_command();
$trace68$;
send_command(command);
$trace69$;
switch_command();
$trace70$;
string$s0;
string$s0;
string$s0;
string$s0;
string$s0;
string$s0;
string$s0;
string$s0;
string$s0;
string$s0;
string$s0;
string$s0;
string$s0;
string$s0;
string$s0;
string$s0;
string$s0;
string$s0;
string$s0;
$s0=str(num);
$s0=$s0;
UNDEFINE($s0);
$s0=$s0;
UNDEFINE($s0);
$s0=$s0;
UNDEFINE($s0);
$s0=$s0;
UNDEFINE($s0);
$s0=$s0;
UNDEFINE($s0);
$s0=$s0;
UNDEFINE($s0);
$s0=$s0;
UNDEFINE($s0);
$s0=$s0;
UNDEFINE($s0);
$s0=$s0;
UNDEFINE($s0);
$s0=$s0;
UNDEFINE($s0);
$s0=$s0;
UNDEFINE($s0);
$s0=$s0;
UNDEFINE($s0);
$s0=$s0;
UNDEFINE($s0);
$s0=$s0;
UNDEFINE($s0);
$s0=$s0;
UNDEFINE($s0);
$s0=$s0;
UNDEFINE($s0);
$s0=$s0;
UNDEFINE($s0);
$s0=$s0;
UNDEFINE($s0);
snum=$s0;
UNDEFINE($s0);
$trace71$;
buf=(("#for"+snum)+":print(\"\")");
$trace72$;
send_command(buf);
$trace73$;
buf="if(True){print(\"\")";
$trace74$;
send_command(buf);
$trace75$;
buf=(("if("+command)+"){print(\"\")");
$trace76$;
send_command(buf);
$trace77$;
switch_command();
$trace78$;
counter=block_end();
$trace79$;
int$ind0;
int$ind0;
int$ind0;
int$ind0;
int$ind0;
int$ind0;
int$ind0;
int$ind0;
int$ind0;
int$ind0;
int$ind0;
int$ind0;
int$ind0;
int$ind0;
int$ind0;
int$ind0;
int$ind0;
int$ind0;
int$ind0;
$ind0=index(command,"{");
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
$ind0=$ind0;
UNDEFINE($ind0);
pos=$ind0;
UNDEFINE($ind0);
$trace80$;
if(-1==pos){println("for: ERROR");
$trace81$;
exit(1);
$trace82$;
};
$trace83$;
pos=(pos-1);
$trace84$;
inc=command[0:pos];
$trace85$;
#next_internal:switch_command();
$trace86$;
if(COMMAND_COUNTER<counter){if(is_for(command)){was_internal_for=True;
$trace87$;
send_command(command);
$trace88$;
switch_command();
$trace89$;
send_command(command);
$trace90$;
switch_command();
$trace91$;
send_command(command);
$trace92$;
switch_command();
$trace93$;
send_command(command);
$trace94$;
internal_counter=block_end();
$trace95$;
switch_command();
$trace96$;
};
$trace97$;
if(was_internal_for){if(COMMAND_COUNTER>internal_counter){was_internal_for=False;
$trace98$;
};
$trace99$;
};
$trace100$;
if(("break"==command)AND(NOT(was_internal_for))){command=(("goto(#for"+snum)+"_end)");
$trace101$;
};
$trace102$;
if(("continue"==command)AND(NOT(was_internal_for))){send_command(inc);
$trace103$;
command=(("goto(#for"+snum)+")");
$trace104$;
};
$trace105$;
send_command(command);
$trace106$;
goto(#next_internal);
$trace107$;
};
$trace108$;
send_command(inc);
$trace109$;
buf=(("goto(#for"+snum)+")");
$trace110$;
send_command(buf);
$trace111$;
send_command(command);
$trace112$;
command=(("#for"+snum)+"_end:print(\"\")");
$trace113$;
send_command(command);
$trace114$;
command="}";
$trace115$;
send_command(command);
$trace116$;
num=(num+1);
$trace117$;
}else{send_command(command);
$trace118$;
};
$trace119$;
goto(#next);
$trace120$;
};
$trace121$;
if(was_for){was_for=False;
$trace122$;
switch_files();
$trace123$;
COMMAND_COUNTER=0;
$trace124$;
goto(#next);
$trace125$;
};
$trace126$;
clear_files();
$trace127$;
};
$trace128$;
main();
