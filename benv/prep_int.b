#import "stdlib/core.b"

string root_source;
string command;

void init(){
	if (exists("benv/import_program.b")){
		root_source = "benv/import_program.b";
	}else{
		root_source = "benv/trace_program.b";	
	};
	SET_SOURCE(root_source);
	SET_DEST("benv/int_program.b");	
};

void finish(){
	UNSET_SOURCE();
	UNSET_DEST();
};

void copy(string source, string dest){
	SET_SOURCE(source);
	SET_DEST(dest);
	string command;
	next_command(command); 
	while (NOT("end" == command)){
		send_command(command);
		next_command(command);	
	};

	UNSET_SOURCE();
	UNSET_DEST();
};

bool is_int(string command){
	stack s;
	stack el;
	stack res;
	stack null; 
	int pos;
	int epos;
	string find;

	find = "int(";
	s = ops(command, find);
	string buf;
	s.pop(buf);
	
	while (NOT("end" == buf)){
		pos = int(buf);
		if (NOT(0 == pos)){
			if (NOT((is_letter(command[(pos - 1)])) OR (is_digit(command[(pos - 1)])))){
				return True;	
			};	
		};
		s.pop(buf);
	};

	return False;
};

void modify(){
	int bpos;
	int epos;
	int number;
	string snumber;
	string sub_command; 

	next_command(command);
	
	while (NOT("end" == command)){
		number = 0;

		if (is_int(command)){
			send_command("print(\"\")");
		};
		
		send_command(command); 
		next_command(command);
	};
};

void main(){
	init();
	modify();
	finish();

	if ("benv/import_program.b" == root_source){
		copy("benv/int_program.b", "benv/import_program.b");
	}else{
		copy("benv/int_program.b", "benv/trace_program.b");
	};
	
	DEL_DEST("benv/int_program.b");

	for (int number; number = 0; exists((("benv/trace/trace_program" + str(number)) + ".b")); number = (number + 1)){
		SET_SOURCE((("benv/trace/trace_program" + str(number)) + ".b"));
		SET_DEST((("benv/trace/int_program" + str(number)) + ".b"));
		modify();
		finish();
		copy((("benv/trace/int_program" + str(number)) + ".b"), (("benv/trace/trace_program" + str(number)) + ".b"));
		DEL_DEST((("benv/trace/int_program" + str(number)) + ".b"));
	};
};

main();
