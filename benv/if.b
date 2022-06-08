#import "stdlib/core.b"

string root_source;
string command;
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
	int o_sum;
	int c_sum;
	int command_len;
	stack obraces;
	stack cbraces;
	string buf;
	string spos;
	int counter;
	int pos; 

	pos = 0;
	counter = COMMAND_COUNTER;
	code = "";

	#block_s:
	pos = 0;
	
	println(command);
	buf = str(o_sum);
	println(buf);
	buf = str(c_sum);
	println(buf);
	
	o_sum = 1;
	c_sum = 0;
	code = (code + command);
	command_len = len(code);
	counter = (counter + 1);

	obraces = ops(code, "{");
	cbraces = ops(code, "}");
	
	#braces_s:
	[print(""), (pos < command_len), goto(#braces_e)];
	spos = str(pos);
	[print(""), in_stack(obraces, spos), goto(#o_plus_end)];
	o_sum = (o_sum + 1);
	#o_plus_end:
	[print(""), in_stack(cbraces, spos), goto(#c_plus_end)];
	c_sum = (c_sum + 1);
	#c_plus_end:
	[goto(#block_e), (o_sum == c_sum), print("")];
	pos = (pos + 1);
	goto(#braces_s);
	#braces_e:
	next_command(command);
	goto(#block_s);
	#block_e:
	SET_COMMAND_COUNTER(COMMAND_COUNTER);
	return counter; 
};

void main(){
	string buf;
	int counter;
	COMMAND_COUNTER = 0;

	init();
	
	next_command(command);
	COMMAND_COUNTER = (COMMAND_COUNTER + 1);
	#main_s:
	[goto(#main_e), ("end" == command), print("")];
	[print(""), (is_if(command)), goto(#next)];
	println("**************************");
	println(command);
	println(get_cond(command));	
	println("---------------------------");
	counter = block_end();
	#next:
	send_command(command);
	next_command(command);
	COMMAND_COUNTER = (COMMAND_COUNTER + 1);
	goto(#main_s);
	#main_e:	

	finish();
};

main();
