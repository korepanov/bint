#import "stdlib/core.b"

string root_source;
string command;
bool first_file;
int COMMAND_COUNTER;
int num;
int exit_num;

void init(){
	root_source = "benv/import_program.b";
	SET_SOURCE(root_source);
	SET_DEST("benv/if_program.b");
};

void finish(){
	UNSET_SOURCE();
	UNSET_DEST();
};

bool is_if(string command){
	string op;
	stack s;
	string buf;

	op = "if(";	
	s = ops(command, op);
	s.pop(buf);
	[goto(#end_true), ("end" == buf), print("")];
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

void SET_COMMAND_COUNTER(int counter){
	int i;
	i = 0;
	string command;
	RESET_SOURCE();
	COMMAND_COUNTER = counter;

	#set_start:
	[print(""), (i < counter), goto(#set_end)];
	next_command(command);
	i = (i + 1);
	goto(#set_start);
	#set_end:
	print("");
};

string get_command(int counter){
	int i;
	string buf;
	string command;

	i = 0;
	RESET_SOURCE();

	#get_command_s:
	[print(""), (i < counter), goto(#get_command_e)];
	next_command(command);
	i = (i + 1);
	goto(#get_command_s);
	#get_command_e:
	SET_COMMAND_COUNTER(COMMAND_COUNTER);
	
	return command;
};

int stack_len(stack s){
	int res;
	string buf; 
	res = 0; 
	
	s.pop(buf);
	#stack_len_s:
	[goto(#stack_len_e), ("end" == buf), print("")];
	res = (res + 1);
	s.pop(buf);
	goto(#stack_len_s);
	#stack_len_e:
	return res; 
};

int block_end(){
	string op1;
	string op2;
	string code;
	string command_buf;
	int o_sum;
	int c_sum;
	int command_len;
	stack obraces;
	stack cbraces;
	string buf;
	string spos;
	int counter;
	int buf_counter;
	int pos; 

	counter = COMMAND_COUNTER;
	command_len = len(command);
	command = command[1:command_len];
	code = command;
	next_command(command);
	counter = (counter + 1);
	code = (code + command[0]);
	o_sum = 1;
	c_sum = 0;
	
	#block_s:
	[goto(#block_e), (o_sum == c_sum), print("")];
	obraces = ops(code, "{");
	cbraces = ops(code, "}");
	
	o_sum = stack_len(obraces);
	c_sum = stack_len(cbraces);
	
	command_buf = command[1:command_len];
	code = (code + command_buf);
	next_command(command); 
	counter = (counter + 1);
	command_buf = command[0]; 
	code = (code + command_buf);
	
	goto(#block_s);
	#block_e:
	buf_counter = (COMMAND_COUNTER - 1);
	SET_COMMAND_COUNTER(buf_counter);
	next_command(command);
	COMMAND_COUNTER = (COMMAND_COUNTER + 1);
	return (counter - 1); 
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

void switch_command(){
	COMMAND_COUNTER = (COMMAND_COUNTER + 1);
	next_command(command);
};

void replace_if(string cond, int stop_pos){
	string buf;
	string snum;
	
	snum = str(num);
	buf = (((("[print(\"\"), " + cond) + ", goto(#_cond") + snum) + "_end)]");
	send_command(buf);
	switch_command(); 

	#replace_clear_if_s:
	[goto(#replace_clear_if_e), ("end" == command), print("")];
	[print(""), (stop_pos == COMMAND_COUNTER), goto(#add_replace_clear_if_mark)];
	buf = (("#_cond" + snum) + "_end:print(\"\")");
	send_command(buf);
	switch_command();
	#add_replace_clear_if_mark:
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
	buf = ((("#_cond_exit") + sexit_num) + ":print(\"\")");
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
	buf = ((("#_cond_exit") + sexit_num) + ":print(\"\")");
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
	DEL_DEST(root_source);
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
