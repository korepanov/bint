#import "stdlib/core.b"

string root_source;
string command;
bool first_file;
int num;
int exit_num;
int br_closed;
int br_opened;

void init(){
	br_closed = 0;
	br_opened = 0;
	get_root_source(root_source);
	SET_SOURCE(root_source);
	SET_DEST("benv/if_program.b");
};

void finish(){
	UNSET_SOURCE();
	UNSET_DEST();
};

bool is_var_def(string command){
	stack s;
	string buf;

	s = reg_find("^(?:(int|float|bool|stack|string)[^\(]*)", command);
	s.pop(buf);
	
	return (NOT("end" == buf)); 
};

string Type(string command){
	stack s;
	string buf;

	s = reg_find("(?:(^int))", command);
	s.pop(buf);
	
	[goto(#int_end), ("end" == buf), print("")];
	return "int";	
	#int_end:

	s = reg_find("(?:(^float))", command);
	s.pop(buf);
	[goto(#float_end), ("end" == buf), print("")];
	return "float";	
	#float_end:

	s = reg_find("(?:(^bool))", command);
	s.pop(buf);
	[goto(#bool_end), ("end" == buf), print("")];
	return "bool";	
	#bool_end:

	s = reg_find("(?:(^stack))", command);
	s.pop(buf);
	[goto(#stack_end), ("end" == buf), print("")];
	return "stack";	
	#stack_end:

	s = reg_find("(?:(^string))", command);
	s.pop(buf);
	[goto(#string_end), ("end" == buf), print("")];
	return "string";	
	#string_end:

	s = reg_find("(?:(^void))", command);
	s.pop(buf);
	[goto(#void_end), ("end" == buf), print("")];
	return "void";	
	#void_end:
	
	print(command);
	print("\n");
	print("Type: ERROR\n");
	exit(1);
};

void check_br(string command){
	string symbol;
	stack s;
	string buf;

	symbol = "{";
	s = ops(command, symbol);
	s.pop(buf);
	
	[goto(#br_opened_end), ("end" == buf), print("")];
	br_opened = (br_opened + 1);
	#br_opened_end:

	symbol = "}";
	s = ops(command, symbol);
	s.pop(buf);
	
	[goto(#br_closed_end), ("end" == buf), print("")];
	br_closed = (br_closed + 1);
	
	#br_closed_end:
	print("");
};

void reset_br(){
	br_opened = 0;
	br_closed = 0;
};

bool is_if(string command){
	string op;
	stack s;
	string buf;
	int ibuf;
	int start_pos;

	op = "if(";	
	s = ops(command, op);
	s.pop(buf);
	[goto(#end_true), ("end" == buf), print("")];
	ibuf = int(buf);
	[goto(#start_true), (0 == ibuf), print("")];
	start_pos = (ibuf - 4);
	[goto(#end_true), (start_pos < 0), print("")];
	buf = command[start_pos:ibuf];
	[print(""), ("else" == buf), goto(#end_true)];
	#start_true:
	return True;
	#end_true:
	return False;
	print("");
};

string get_cond(string command){
	string op;
	stack s;
	string buf;
	int start_pos;
	int end_pos; 

	op = "if(";
	s = ops(command, op);
	s.pop(buf);
	start_pos = int(buf);
	start_pos = (start_pos + 2);
	end_pos = func_end(command, start_pos);
	end_pos = (end_pos + 1);
	buf = command[start_pos:end_pos];

	return buf;
};

string if_type(string command){
	int command_len;
	string prefix; 
	
	command_len = len(command);
	[print(""), (1 == command_len), goto(#not_clear)];
	
	return "clear";
	#not_clear:
	prefix = command[1:7];
	
	[print(""), ("elseif" == prefix), goto(#not_elseif)];
	return "elseif";

	#not_elseif:
	prefix = command[1:5];
	[print(""), ("else" == prefix), goto(#if_type_error)];
	return "else";
	
	#if_type_error:
	return "error";	
};

void switch_files(){
	finish();
	[print(""), (first_file), goto(#first_end)];
	SET_SOURCE("benv/if_program.b");
	SET_DEST("benv/if_program2.b");
	first_file = False;
	goto(#switch_files_e);
	#first_end:
	SET_SOURCE("benv/if_program2.b");
	SET_DEST("benv/if_program.b");
	first_file = True; 
	#switch_files_e:
	print("");
};

void replace_if(string cond, int stop_pos){
	string buf;
	string snum;
	stack args_to_undefine;
	string arg_type;
	int command_len;
	string arg_name;
	
	snum = str(num);
	buf = (((("[print(\"\"), " + cond) + ", goto(#_cond") + snum) + "_end)]");
	send_command(buf);
	switch_command(); 

	#replace_clear_if_s:
	[goto(#replace_clear_if_e), ("end" == command), print("")];
	[print(""), (stop_pos == COMMAND_COUNTER), goto(#add_replace_clear_if_mark)];
	
	args_to_undefine.pop(arg_name);
	#un:
	[goto(#un_end), ("end" == arg_name), print("")];
	buf = (("UNDEFINE(" + arg_name) + ")");
	send_command(buf);
	args_to_undefine.pop(arg_name); 
	goto(#un);
	#un_end:
	
	buf = (("#_cond" + snum) + "_end:print(\"\")");
	send_command(buf);
	switch_command();
	#add_replace_clear_if_mark:
	check_br(command);

	[print(""), ((is_var_def(command))AND(br_closed == br_opened)), goto(#pop_e)];
	arg_type = Type(command);
	int type_len;
	type_len = len(arg_type);
	command_len = len(command);
	arg_name = command[type_len:command_len];
	args_to_undefine.push(arg_name);
	#pop_e:

	send_command(command);
	switch_command();
	goto(#replace_clear_if_s);
	#replace_clear_if_e:
	COMMAND_COUNTER = 0;
	switch_files();
};

void replace_elseif(string cond, int stop_pos){
	string buf;
	int ibuf;
	string snum; 
	string t;
	string bcommand;
	int this_stop_pos;
	string sexit_num;
	
	sexit_num = str(exit_num);
	snum = str(num);
	buf = (((("[print(\"\"), " + cond) + ", goto(#_cond") + snum) + "_end)]");
	send_command(buf);
	switch_command(); 
	
	#replace_elseif_s:
	[goto(#replace_elseif_e), ("end" == command), print("")];
	[print(""), (stop_pos == COMMAND_COUNTER), goto(#add_replace_elseif_mark)];
	sexit_num = str(exit_num);
	buf = (("goto(#_cond_exit" + sexit_num) + ")");
	send_command(buf);
	buf = (("#_cond" + snum) + "_end:print(\"\")");
	send_command(buf);
	num = (num + 1);
	t = if_type(command);
	[print(""), ("elseif" == t), goto(#find_block_e)];
	this_stop_pos = block_end();
	#find_block_e:
	bcommand = command;
	switch_command();
	[goto(#replace_elseif_e), (NOT("elseif" == t)), print("")];
	[print(""), ("elseif" == t), goto(#elif_end)];
	cond = get_cond(bcommand);
	snum = str(num);
	buf = (((("[print(\"\"), " + cond) + ", goto(#_cond") + snum) + "_end)]");
	send_command(buf);
	stop_pos = this_stop_pos; 
	goto(#replace_elseif_s);
	#elif_end:
	print("");
	#add_replace_elseif_mark:
	send_command(command);
	switch_command();
	goto(#replace_elseif_s);
	#replace_elseif_e:
	t = if_type(bcommand);
	[print(""), ("else" == t), goto(#restore_end)];
	ibuf = (COMMAND_COUNTER - 2);
	SET_COMMAND_COUNTER(ibuf);
	switch_command();
	stop_pos = block_end();
	switch_command();
	#restore_end:
	[print(""), ("else" == t), goto(#final_cond_end)];
	#final_cond_s:
	[goto(#final_cond_end), (stop_pos == COMMAND_COUNTER), print("")];
	send_command(command);
	switch_command();
	goto(#final_cond_s);
	#final_cond_end:
	[print(""), ("else" == t), goto(#else_end)];
	switch_command();
	#else_end:
	sexit_num = str(exit_num);
	buf = (("#_cond_exit" + sexit_num) + ":print(\"\")");
	exit_num = (exit_num + 1);
	send_command(buf);
	#ts:
	[goto(#te), ("end" == command), print("")];
	send_command(command);
	switch_command();
	goto(#ts);
	#te:
	COMMAND_COUNTER = 0;
	switch_files();
};

void replace_else(string cond, int stop_pos){
	string buf;
	string snum;
	string sexit_num;
	int ibuf; 
	int pos; 
	
	pos = -1;
	snum = str(num);
	buf = (((("[print(\"\"), " + cond) + ", goto(#_cond") + snum) + "_end)]");
	send_command(buf);
	switch_command(); 

	#replace_else_s:
	[goto(#replace_else_e), ("end" == command), print("")];
	[print(""), (stop_pos == COMMAND_COUNTER), goto(#add_replace_else_mark)];
	sexit_num = str(exit_num);
	buf = (("goto(#_cond_exit" + sexit_num) + ")");
	send_command(buf); 
	buf = (("#_cond" + snum) + "_end:print(\"\")");
	send_command(buf);
	switch_command();

	ibuf = (COMMAND_COUNTER - 2);
	SET_COMMAND_COUNTER(ibuf);
	switch_command();
	pos = block_end();
	switch_command();
	
	#add_replace_else_mark:
	[print(""), (pos == COMMAND_COUNTER), goto(#figure_brace_end)];
	sexit_num = str(exit_num);
	buf = (("#_cond_exit" + sexit_num) + ":print(\"\")");
	exit_num = (exit_num + 1);
	send_command(buf);
	switch_command();
	#ets:
	[goto(#ete), ("end" == command), print("")];
	send_command(command);
	switch_command();
	goto(#ets);
	#ete:
	goto(#replace_else_e);
	#figure_brace_end:
	send_command(command);
	switch_command();
	goto(#replace_else_s);
	#replace_else_e:
	COMMAND_COUNTER = 0;
	switch_files();
};

void clear_files(){
	[goto(#first_file), (first_file), print("")];
	switch_files();
	switch_command();
	#clear_files_s:
	[goto(#clear_files_e), ("end" == command), print("")];
	send_command(command);
	switch_command();
	goto(#clear_files_s);

	#first_file:
	print("");
	#clear_files_e:
	finish();
	DEL_DEST("benv/if_program2.b");
};

void main(){
	string buf;
	string cond;
	int counter;
	int buf_counter;
	string snum;
	string t;
	bool is_stop;
	
	first_file = True;
	init();
	num = 0;
	exit_num = 0;
	COMMAND_COUNTER = 0;

	#again_s:
	is_stop = True;
	
	
	switch_command(command);
	#main_s:
	[goto(#main_e), ("end" == command), print("")];
	[print(""), (is_if(command)), goto(#next)];
	is_stop = False;
	cond = get_cond(command);
	counter = block_end();
	buf = get_command(counter);
	t = if_type(buf);	
	[print(""), ("error" == t), goto(#error_end)];
	println("ERROR in the if type\n");
	goto(#total_e);
	#error_end:
	[print(""), ("clear" == t), goto(#if_end)];
	replace_if(cond, counter);
	num = (num + 1);
	goto(#main_e);
	#if_end:
	[print(""), ("elseif" == t), goto(#elseif_end)];
	replace_elseif(cond, counter);
	goto(#main_e);
	#elseif_end:
	replace_else(cond, counter);
	num = (num + 1);
	goto(#main_e); 
	#next:
	send_command(command);
	switch_command();
	goto(#main_s);	
	#main_e:
	[print(""), (is_stop), goto(#again_s)];
	#total_e:
	clear_files();
};

main();
