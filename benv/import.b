string root_source;
string first_command;
string file;
string translate_mode;
stack gimports;
stack pimports;

void init(){
	string symbol;
	int pos; 

	get_root_source(root_source);
	
	pos = len(root_source);
	pos = (pos - 1);
	symbol = root_source[pos];

	if ("d" == symbol){
		translate_mode = "debug";
		root_source = root_source[0:pos];	
	}else if ("v" == symbol){
		translate_mode = "validate";
		root_source = root_source[0:pos];
	}else{
		translate_mode = "release";
	};

	SET_SOURCE(root_source);
		
};

void finish(){
	UNSET_DEST();
};

bool in_stack(stack s, string el){
	string buf;
	bool res;

	res = False;
	s.pop(buf);
	#in_stack_s:
	[goto(#in_stack_e), ("end" == buf), print("")];
	[print(""), (el == buf), goto(#no)];
	res = True;
	goto(#in_stack_e);
	#no:
	s.pop(buf);
	goto(#in_stack_s);	
	#in_stack_e:
	return res;
};

int find_import_end(string command){
	int command_len;
	int res;

	command_len = len(command);	
	
	command = command[8:command_len];
	res = index(command, "\"");
	res = (res + 8);	
 
	return res;
};

stack reverse(stack s){
	string buf;
	stack res; 
	
	s.pop(buf);
	#reverse_s:
	[goto(#reverse_e), ("end" == buf), print("")];
	res.push(buf);
	s.pop(buf);
	goto(#reverse_s);
	#reverse_e:
	return res;
};

stack get_imports(){
	stack res; 
	int number;
	string import_name;
	string command;
	int command_len;
	int import_end;

	next_command(command);
	#get_imports_s:
	[goto(#get_imports_e), ("end" == command), print("")];
	number = index(command, "#import");
	[print(""), (0 == number), goto(#get_imports_e)];
	command_len = len(command);
	import_end = find_import_end(command);
	import_name = command[8:import_end];
	res.push(import_name);
	import_end = (import_end + 1);
	command = command[import_end:command_len];
	goto(#get_imports_s);
	#get_imports_e:
	RESET_SOURCE();
	first_command = command;
	res = reverse(res);	
	return res; 
};

void import_union(string file){
	string command;
	stack imports;
	string import;
	string f;
	int number;
	int import_end;
	int command_len;

	imports = get_imports();
	imports.pop(import);
	
	do{
		while ((in_stack(gimports, import))AND(NOT("end" == import))){
			imports.pop(import);
		};
		if ("end" == import){
			UNSET_SOURCE();
		}else{
			gimports.push(import);
			UNSET_SOURCE();
			SET_SOURCE(import);
			import_union(import);
			
			if (("debug" == translate_mode) OR ("validate" == translate_mode)){
				f = (("$file " + import) + "$");
				send_command(f);
			};
			SET_SOURCE(import);
			next_command(command);
			
			while (NOT("end" == command)){
				number = index(command, "#import");
				while (0 == number){
					import_end = find_import_end(command);
					import_end = (import_end + 1);
					command_len = len(command);
					command = command[import_end:command_len];
					number = index(command, "#import");
				};
				send_command(command);
				next_command(command);		
			};	
			UNSET_SOURCE();
		};
		imports.pop(import);
		if (NOT("end" == import)){
			SET_SOURCE(import);
		};
	}while(NOT("end"==import));
	
};

void file_union(string file){
	string command; 
	int number; 
	int import_end;
	int command_len; 
	string f; 

	import_union(file);
	SET_SOURCE(file);
	if (("debug" == translate_mode) OR ("validate" == translate_mode)){
		f = (("$file " + file) + "$");
		send_command(f);
	};
	next_command(command);
	while (NOT("end" == command)){
		number = index(command, "#import");
		while (0 == number){
			import_end = find_import_end(command);
			import_end = (import_end + 1);
			command_len = len(command);
			command = command[import_end:command_len];
			number = index(command, "#import");
		};
		send_command(command);
		next_command(command);
	}; 
	UNSET_SOURCE();	
};

void get_file_content(string file){
	string import;
	stack imports;
	string command;
	int number;
	
	imports = get_imports();
	imports.pop(import);
	UNSET_SOURCE();
	
	while(NOT("end"==import)){
		SET_SOURCE(import);
		f = (("$file " + import) + "$");
		send_command(f);
		next_command(command);
		while (NOT("end" == command)){
			number = index(command, "#import");
			while (0 == number){
				import_end = find_import_end(command);
				import_end = (import_end + 1);
				command_len = len(command);
				command = command[import_end:command_len];
				number = index(command, "#import");
			};
			send_command(command);
			next_command(command);
		}; 
		imports.pop(import);
	};
	
	SET_SOURCE(file);
	f = (("$file " + file) + "$");
	send_command(f);
	next_command(command);
	
	while (NOT("end" == command)){
		number = index(command, "#import");
		while (0 == number){
			import_end = find_import_end(command);
			import_end = (import_end + 1);
			command_len = len(command);
			command = command[import_end:command_len];
			number = index(command, "#import");
		};
		send_command(command);
		next_command(command);
	}; 
};

void make_parts(string file, int number){
	string command;
	stack imports;
	string import;
	string f;
	int import_end;
	int command_len;
	int file_number;
	string snumber;
	bool got_content;
	string dest; 
	bool cond;
	string temp;

	imports = get_imports();
	imports.pop(import);
	
	do{
		if ("end" == import){
			UNSET_SOURCE();
		}else{	
			if (NOT(got_content)){
				got_content = True;
				snumber = str(number);
				dest = (("benv/import/import_program" + snumber) + ".b");
				SET_DEST(dest);
				number = (number + 1);
				get_file_content(file);
				UNSET_DEST();
			};
			UNSET_SOURCE();
			SET_SOURCE(import);
			cond = NOT(in_stack(pimports, import));
			temp = str(cond);
			print(temp);
			print("\n");
			if (cond){
				pimpots.push(import);
				make_parts(import, number);
			}else{
				print("import cycle not allowed\n");
				print("import in ");
				print(file);
				print(": ");
				print(import);
				exit(1);
			};
			
			SET_SOURCE(import);
			UNSET_SOURCE();
		};
		imports.pop(import);
		if (NOT("end" == import)){
			SET_SOURCE(import);
		};
	}while(NOT("end"==import));
};

void del_files(){
	int number;
	string snumber;
	bool e;
	string dest; 
	
	snumber = str(number);
	dest = (("benv/import/import_program" + snumber) + ".b");
	e = exists(dest);
	
	while(e){
		DEL_DEST(dest);
		number = (number + 1);
		snumber = str(number);
		dest = (("benv/import/import_program" + snumber) + ".b");
		e = exists(dest);
	};
	
};

void main(){
	init();
	del_files();
	
	if ("validate" == translate_mode){
		make_parts(root_source, 0);
		SET_SOURCE(root_source);
	};
	
	SET_DEST("benv/import_program.b");
	file_union(root_source);
	finish();	
};

main();
