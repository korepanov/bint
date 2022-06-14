#import "stdlib/core.b"

string root_source;
string command;
bool first_file;
int COMMAND_COUNTER;

void init(){
	get_root_source(root_source);
	SET_SOURCE(root_source);
	SET_DEST("benv/program.basm");
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
	SET_SOURCE("benv/program.basm");
	SET_DEST("benv/program2.basm");
	first_file = False;
	goto(#switch_files_e);
	#first_end:
	SET_SOURCE("benv/program2.basm");
	SET_DEST("benv/program.basm");
	first_file = True; 
	#switch_files_e:
	print("");
};

void main(){
	string buf;
	string cond;
	int counter;
	int buf_counter;
	int num;
	string snum;
	string t;
	bool is_stop;
	
	first_file = True;
	init();
	#again_s:
	COMMAND_COUNTER = 0;
	num = 0;
	is_stop = True;
	
	
	next_command(command);
	COMMAND_COUNTER = (COMMAND_COUNTER + 1);
	#main_s:
	[goto(#main_e), ("end" == command), print("")];
	[print(""), (is_if(command)), goto(#next)];
	is_stop = False;
	cond = get_cond(command);
	counter = block_end();
	buf = get_command(counter);
	t = if_type(buf);
	[print(""), (("clear" == t) OR ("elseif" == t)), goto(#replace_clear_end)];
	snum = str(num);
	command = (((("[print(\"\"), " + cond) + ", goto(#cond") + snum) + "_end)]");
	send_command(command);
	buf_counter = (counter - 1);
	#next_block_command_s:
	[print(""), (COMMAND_COUNTER < buf_counter), goto(#next_block_command_e)];
	next_command(command);
	send_command(command);
	COMMAND_COUNTER = (COMMAND_COUNTER + 1);
	goto(#next_block_command_s);
	#next_block_command_e:
	command = (("#cond" + snum) + "_end:print(\"\")");
	send_command(command);
	next_command(command);
	COMMAND_COUNTER = (COMMAND_COUNTER + 1);
	[print(""), ("elseif" == t), goto(#elseif_end)];
	
	goto(#replace_clear_end);
	#elseif_end:
	next_command(command);
	COMMAND_COUNTER = (COMMAND_COUNTER + 1);
	#replace_clear_end:	
	[print(""), ("else" == t), goto(#replace_else_end)];
	snum = str(num);
	command = (((("[print(\"\"), " + cond) + ", goto(#cond") + snum) + "_end)]");
	send_command(command);
	buf_counter = (counter - 1);
	#next_block_command2_s:
	[print(""), (COMMAND_COUNTER < buf_counter), goto(#next_block_command2_e)];
	next_command(command);
	send_command(command);
	COMMAND_COUNTER = (COMMAND_COUNTER + 1);
	goto(#next_block_command2_s);
	#next_block_command2_e:
	command = (((("#cond" + snum) + "_end:goto(#other") + snum) + "_end)");
	send_command(command);
 	command = (("#other" + snum) + ":print(\"\")");
	send_command(command);
	counter = block_end();
	buf_counter = (counter - 1);
	#next_block_command3_s:
	[print(""), (COMMAND_COUNTER < buf_counter), goto(#next_block_command3_e)];
	next_command(command);
	send_command(command);
	COMMAND_COUNTER = (COMMAND_COUNTER + 1);
	goto(#next_block_command3_s);
	#next_block_command3_e:
	command = (("#other" + snum) + "_end:print(\"\")");
	send_command(command);
	next_command(command);
	next_command(command);
	#replace_else_end:
	num = (num + 1); 
	#next:
	send_command(command);
	next_command(command);
	COMMAND_COUNTER = (COMMAND_COUNTER + 1);
	goto(#main_s);
	#main_e:
	[goto(#again_e), (is_stop), print("")];
	switch_files();
	goto(#again_s);	
	#again_e:
	finish();
};

main();
