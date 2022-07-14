void init(){
	string root_source;
	get_root_source(root_source);
	SET_SOURCE(root_source);
	SET_DEST("benv/recurs_program.b");
};

void finish(){
	UNSET_SOURCE();
	UNSET_DEST();
};

void main(){
	init();
	string command;
	next_command(command);
	
	#main_s:
	print("");
	if (NOT("end" == command)){
		send_command(command);
		next_command(command);
		goto(#main_s);
	};
	finish();
};

main();
