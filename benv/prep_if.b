#import "stdlib/core.b"

string root_source;
int COMMAND_COUNTER;
string command;

void init(){
	COMMAND_COUNTER = 0;
	root_source = "benv/prep_func_program.b";
	SET_SOURCE(root_source);
	SET_DEST("benv/prep_if_program.b");
};

void finish(){
	UNSET_SOURCE();
	UNSET_DEST();
	DEL_DEST(root_source);
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

void switch_command(){
	COMMAND_COUNTER = (COMMAND_COUNTER + 1);
	next_command(command);
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
	
	command_len = len(command);
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

bool is_else(string command){
	string op;
	stack s;
	string buf;
	int ibuf;
	int start_pos;

	op = "}else{";
	s = ops(command, op); 
	s.pop(buf);
	[goto(#end_else_true), ("end" == buf), print("")];
	ibuf = int(buf);
	[print(""), (0 == ibuf), goto(#end_else_true)];
	return True;
	#end_else_true:
	return False;
};

void send_new_command(string command){
	string op;
	string buf;
	string command_buf; 
	stack s; 
	int pos;
	string new_command;
	int command_len;
	op = "{";

	#send_new_command_s:
	s = ops(command, op);
	s.pop(buf);
	
	[goto(#send_new_command_e), ("end" == buf), print("")];
	
	pos = int(buf);
	pos = (pos + 1);
	command_buf = command[0:pos];
	new_command = (command_buf + "print(\"\")");
	send_command(new_command);

	command_len = len(command);
	command_buf = command[pos:command_len];
	command = command_buf;
	goto(#send_new_command_s);
	#send_new_command_e:
	send_command(command);
	print("");
};

string var_type(string command){
	string buf;
	stack s;
	
	s = reg_find("(?:^int)", command);
	s.pop(buf);
	[goto(#int_type_end), ("end" == buf), print("")];
	return "int";
	#int_type_end:

	s = reg_find("(?:^float)", command);
	s.pop(buf);
	[goto(#float_type_end), ("end" == buf), print("")];
	return "float";
	#float_type_end:

	s = reg_find("(?:^bool)", command);
	s.pop(buf);
	[goto(#bool_type_end), ("end" == buf), print("")];
	return "bool";
	#bool_type_end:

	s = reg_find("(?:^stack)", command);
	s.pop(buf);
	[goto(#stack_type_end), ("end" == buf), print("")];
	return "stack";
	#stack_type_end:

	s = reg_find("(?:^string)", command);
	s.pop(buf);
	[goto(#string_type_end), ("end" == buf), print("")];
	return "string";
	#string_type_end:
	return "";
};

void undefine_vars(){
	int old_COMMAND_COUNTER;
	int counter;
	int internal_counter;
	int pos;
	string op;
	string buf; 
	string this_command;
	int command_len;
	stack s;
	string T;
	string var;
	stack vars;

	this_command = command;
	old_COMMAND_COUNTER = COMMAND_COUNTER;
	
	counter = block_end();
	
	op = "{";

	#undefine_vars_loop_s:
	s = ops(this_command, op);
	s.pop(buf);
	[goto(#undefine_vars_loop_e), ("end" == buf), print("")];
	pos = int(buf);
	pos = (pos + 1);
	command_len = len(command);
	this_command = command[pos:command_len];
	goto(#undefine_vars_loop_s);
	#undefine_vars_loop_e:
	T = var_type(this_command);
	pos = len(T);
	command_len = len(this_command);
	[print(""), (pos > 0), goto(#not_var)];
	var = this_command[pos:command_len];
	vars.push(var);	
	#not_var:
	print("");
	#vars_loop_s:
	[print(""), (COMMAND_COUNTER < counter), goto(#vars_loop_e)];
	switch_command();
	[print(""), ((is_if(command)) OR (is_else(command))), goto(#not_internal)];
	internal_counter = block_end();
	SET_COMMAND_COUNTER(internal_counter);
	#not_internal:
	T = var_type(command);
	pos = len(T);
	command_len = len(command);
	[print(""), (pos > 0), goto(#not_var2)];
	var = command[pos:command_len];
	vars.push(var);	
	#not_var2:
	goto(#vars_loop_s);
	#vars_loop_e:
	
	SET_COMMAND_COUNTER(old_COMMAND_COUNTER);
	switch_command();

	#modify_s:
	[print(""), (COMMAND_COUNTER < counter), goto(#modify_e)];
	send_command(command);
	switch_command();	
	goto(#modify_s);
	#modify_e:

	vars.pop(var);
	#get_var_s:
	[goto(#get_var_e), ("end" == var), print("")];
	this_command = (("UNDEFINE(" + var) + ")");
	send_command(this_command);
	vars.pop(var);	
	goto(#get_var_s);
	#get_var_e:
	print("");
	
};

void main(){
	int counter; 
	
	init();
	switch_command();
	#main_s:
	[goto(#main_e), ("end" == command), print("")];
	[print(""), ((is_if(command)) OR (is_else(command))), goto(#not_cond)];
	send_new_command(command);
	undefine_vars();
	goto(#main_s);
	#not_cond:
	send_command(command);
	switch_command();
	goto(#main_s);
	#main_e:
	finish();
};

main();
