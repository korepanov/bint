string root_source;
string first_command;
string file;
string translate_mode;

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
	SET_DEST("benv/import_program.b");

	if ("validate" == translate_mode){
		send_command("$validate$");	
	};

	if ("debug" == translate_mode){
		send_command("$debug$");	
	};
		
};

void finish(){
	UNSET_DEST();
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
	return reverse(res);	
};

void file_union(string file){
	string command;
	stack imports;
	string import;
	string f;
	int number;
	int import_end;
	int command_len;

	imports = get_imports();
	imports.pop(import);
	if ("end" == import){
		UNSET_SOURCE();
	}else{
		UNSET_SOURCE();
		SET_SOURCE(import);
		file_union(import);
		
		if (("debug" == translate_mode) OR ("validate" == translate_mode)){
			f = ("$file$ " + import);
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
	
};

void main(){
	init();
	file_union(root_source, False);	
	finish();
};

main();
