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

void main(){
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
