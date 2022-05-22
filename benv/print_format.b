#import "stdlib/core.b"

string root_source;

void init(){
	get_root_source(root_source);
	SET_SOURCE(root_source);
	SET_DEST("benv/print_format_program.b");
};

void finish(){
	UNSET_SOURCE();
	UNSET_DEST();
};

void replace_print(string command){
	stack s;
	string buf;

	s = ops(command, "print");

	s.pop(buf);
	#replace_s:
	[goto(#replace_e), ("end" == buf), print("")];
	println(command);
	println(buf);
	s.pop(buf);	
	goto(#replace_s);
	#replace_e:	

	print("");
};

void main(){
	string command;
	init();
	
	next_command(command);
	#main_s:
	[goto(#main_e), ("end" == command), print("")];
	replace_print(command);
	send_command(command);
	next_command(command);
	goto(#main_s);
	#main_e:
	finish();
};

main();
