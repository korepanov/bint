bool is_release;
string command; 
int COMMAND_COUNTER;

void init(){
	is_release = False;
	COMMAND_COUNTER = 0;

	SET_SOURCE("benv/import_program.b");
	SET_DEST("benv/trace_program.b");
	next_command(command);
	
	if (NOT(("$debug$" == command)OR("$validate$" == command))){
		is_release = True;
	}else{
		send_command(command);
	};

};

void finish(){
	UNSET_SOURCE();
	UNSET_DEST();
	DEL_DEST("benv/import_program.b");
};

void copy(){
	#copy_s:
	print("");
	if(NOT("end" == command)){
		send_command(command);
		next_command(command);
		goto(#copy_s);
	};
};

void switch_command(){
	COMMAND_COUNTER = (COMMAND_COUNTER + 1);
	next_command(command);
};

void modify(){
	string trace; 
	string scommand_counter;

	#modify_s:
	switch_command();
	if ("end" == command){
		goto(#modify_e);	
	};
	
	trace = command[0:6];
	if ("$file$" == trace){
		send_command(command);
		COMMAND_COUNTER = 0;
	}else{
		scommand_counter = str(COMMAND_COUNTER);
		trace = ("$trace$" + scommand_counter);
		send_command(trace);
		send_command(command);
	};
	goto(#modify_s);

	#modify_e:
	print("");
};

void main(){
	init();
	if (is_release){
		copy();
	}else{
		modify();
	};
	finish();
};

main();
