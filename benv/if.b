#import "stdlib/core.b"

string root_source;
string command;

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

void main(){
	string buf;

	init();
	
	next_command(command);
	#main_s:
	[goto(#main_e), ("end" == command), print("")];
	buf = str(is_if(command));
	println(command);
	println(buf);	
	send_command(command);
	next_command(command);
	goto(#main_s);
	#main_e:	

	finish();
};

main();
