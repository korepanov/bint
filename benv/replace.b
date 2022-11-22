#import "stdlib/core.b"

string root_source;
string command;

void init(){
	root_source = "program.b";
	SET_SOURCE(root_source);
	SET_DEST("benv/replace_program.b");	
};

void finish(){
	UNSET_SOURCE();
	UNSET_DEST();
};

void main(){
	init();
	next_command(command);
	send_command(command);
	finish();
};

main();
