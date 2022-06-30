#import "stdlib/core.b"

string root_source;

void init(){
	get_root_source(root_source);
	SET_SOURCE(root_source);
	SET_DEST("benv/prep_if_program.b");
};

void finish(){
	UNSET_SOURCE();
	UNSET_DEST();
};

bool is_cond(string command){
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
	string command;

	init();
	next_command(command);
	#main_s:
	[goto(#main_e), ("end" == command), print("")];
	send_command(command);
	next_command(command);
	goto(#main_s);
	#main_e:
	finish();
};

main();
