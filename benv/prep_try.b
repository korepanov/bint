string root_source;

void init(){
	root_source = "benv/prep_if_program.b";
	SET_SOURCE(root_source);
	SET_DEST("benv/prep_try_program.b");
};

void finish(){
	UNSET_SOURCE();
	UNSET_DEST();
	//DEL_DEST(root_source);
}; 

void main(){
	init();
	
	string command;
	next_command(command);

	while (NOT(command == "end")){
		send_command(command); 
		next_command(command);	
	};
	


	finish(); 
};

main();
