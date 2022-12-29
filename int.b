#import "stdlib/core.b"

string root_source;
string command;

void init(){
	root_source = "benv/import_program.b";
	SET_SOURCE(root_source);
	SET_DEST("benv/int_program.b");	
};

void finish(){
	UNSET_SOURCE();
	UNSET_DEST();
	DEL_DEST(root_source); 
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

stack int_poses(string command){
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
		epos = func_end(command, (pos + 3));
		el.push(epos);
		el.push(pos);
		res.push(el);
		el = null;	
		s.pop(buf);	
	};

	return res;
};

string modify_command(string command, string sub_command, int bpos, int epos){
	string new_command;
	string buf;
	int command_len;
 
	new_command = command[0:bpos];
	new_command = (new_command + sub_command);
	command_len = len(command);
	buf = command[epos:command_len];
	new_command = (new_command + buf);
	
	return new_command;
};

void main(){
	stack s;
	stack el;
	int bpos;
	int epos;
	string buf;
	int number;
	string snumber;
	string sub_command; 

	init();
	next_command(command);
	
	while (NOT("end" == command)){
		number = 0;
		s = int_poses(command);
		s.pop(el);
		el.pop(buf);

		while (NOT("end" == buf)){
			s = int_poses(command);
			s.pop(el);
			el.pop(buf);
			bpos = int(buf);
			el.pop(buf);
			epos = int(buf);
			epos = (epos + 1);
			snumber = str(number);
			buf = ("int $I" + snumber);
			send_command(buf);
			buf = command[bpos:epos];
			buf = ((("$I" + snumber) + "=") + buf);
			send_command(buf);
			sub_command = ("$I" + snumber);
			command = modify_command(command, sub_command, bpos, epos);
			number = (number + 1);
			s.pop(el);
			el.pop(buf);
		};
		send_command(command);

		for (int i; i = 0; i < number; i = (i + 1)){
			string b;
			snumber = str(i);
			b = (("UNDEFINE($I" + snumber) + ")");
			send_command(b);		
		};

		next_command(command);
	};
	finish();
	copy("benv/int_program.b", root_source);
	DEL_DEST("benv/int_program.b"); 
};

main();
