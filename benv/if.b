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

int block_end(){
	string command; 
	string op1;
	string op2;
	int o_sum;
	int c_sum;
	stack obraces;
	stack cbraces;
	string buf; 
	int counter;
	counter = COMMAND_COUNTER;
	#block_s:
	next_command(command);
	println(command);
	counter = (counter + 1);
	o_sum = 1;
	c_sum = 0;

	obrace = "{";
	cbrace = "}";

	obraces = ops(command, obrace);
	cbraces = ops(command, cbrace);
	cbraces.pop(buf);
	println(buf);

	
	obraces.pop(buf);
	
	#obrace_s:
	[goto(#obrace_e), ("end" == buf), print("")];
	o_sum = (o_sum + 1);
	obraces.pop(buf);
	goto(#obrace_s);
	#obrace_e:
	buf = str(o_sum);
	println(buf);

	cbraces.pop(buf);
	println(buf);
	#cbrace_s:
	[goto(#cbrace_e), ("end" == buf), print("")];
	c_sum = (c_sum + 1);
	cbraces.pop(buf);
	goto(#cbrace_s);
	#cbrace_e:
	
	buf = str(c_sum);
	println(buf);
	[goto(#block_e), (o_sum == c_sum), print("")];
	goto(#block_s);
	#block_e:
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
