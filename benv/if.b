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

bool is_else(){
	string op;
	stack s;
	string buf;

	op = "}";
	s = ops(command, op);
	s.pop(buf);

	#is_else_s:
	[print(""), ("end" == buf), goto(#is_else_e)];
	
	#is_else_e:	
	
};

void main(){
	string buf;
	COMMAND_COUNTER = 0;

	init();
	
	next_command(command);
	COMMAND_COUNTER = (COMMAND_COUNTER + 1);
	#main_s:
	[goto(#main_e), ("end" == command), print("")];
	[print(""), (is_if(command)), goto(#next)];
	println(command);
	println(get_cond(command));	
	#next:
	send_command(command);
	next_command(command);
	COMMAND_COUNTER = (COMMAND_COUNTER + 1);
	goto(#main_s);
	#main_e:	

	finish();
};

main();
