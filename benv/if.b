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
	SET_COMMAND_COUNTER(COMMAND_COUNTER);
	return (counter - 1); 
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
	counter = block_end();
	buf = str(COMMAND_COUNTER);
	println(buf);
	buf = str(counter);
	println(buf);
	#next:
	send_command(command);
	next_command(command);
	COMMAND_COUNTER = (COMMAND_COUNTER + 1);
	goto(#main_s);
	#main_e:	

	finish();
};

main();
