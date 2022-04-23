stack reverse(stack s){
	string buf;
	stack res; 
	
	s.pop(buf);
	#_reverse_s:
	[goto(#_reverse_e), ("end" == buf), print("")];
	res.push(buf);
	s.pop(buf);
	goto(#_reverse_s);
	#_reverse_e:
	return res;
};
