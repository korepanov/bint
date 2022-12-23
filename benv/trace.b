bool is_release;
string command; 
int COMMAND_COUNTER;

void init(){
	is_release = False;
	COMMAND_COUNTER = 0;

	SET_SOURCE("benv/import_program.b");
	SET_DEST("benv/trace_program.b");
	next_command(command);
	send_command(command);
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
	int command_len;

	#modify_s:
	trace = "";
	switch_command();
	if ("end" == command){
		goto(#modify_e);	
	};
	
	command_len = len(command);
	
	if (command_len > 5){
		trace = command[0:5];
	};
	if ("$file" == trace){
		send_command(command);
		COMMAND_COUNTER = 0;
	}else{
		scommand_counter = str(COMMAND_COUNTER);
		trace = (("$trace" + scommand_counter) + "$");
		send_command(trace);
		send_command(command);
	};
	goto(#modify_s);

	#modify_e:
	print("");
};

void modify_import(){
	int number;
	string snumber;
	string file;
	string dest;
	bool e;
	
	snumber = str(number);
	file = (("benv/import/import_program" + snumber) + ".b");
	e = exists(file);
	while(e){
		SET_SOURCE(file);
		dest = (("benv/trace/trace_program" + snumber) + ".b");
		SET_DEST(dest);
		modify();
		number = (number + 1);
		snumber = str(number);
		file = (("benv/import/import_program" + snumber) + ".b");
		e = exists(file);
	};
	
};

void del_files(){
	int number;
	string snumber;
	bool e;
	string dest; 
	
	snumber = str(number);
	dest = (("benv/trace/trace_program" + snumber) + ".b");
	e = exists(dest);
	
	while(e){
		DEL_DEST(dest);
		number = (number + 1);
		snumber = str(number);
		dest = (("benv/trace/trace_program" + snumber) + ".b");
		e = exists(dest);
	};
	
};

void main(){
	init();
	if (is_release){
		copy();
	}else{
		modify();
	};
	finish();
	del_files();
	modify_import();
	SET_SOURCE("benv/trace_program.b");
	SET_DEST("benv/import_program.b");
	copy(); 
};

main();
