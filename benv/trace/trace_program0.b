$filestdlib/files.b$;
$trace1$;
intopen_file(stringfile_path,stringmode){intres;
$trace2$;
intm;
$trace3$;
if(mode=="read"){m=0;
$trace4$;
}elseif(mode=="write"){m=1;
$trace5$;
}elseif(mode=="append"){m=2;
$trace6$;
}else{error="no such mode for opening file";
$trace7$;
if($toPanic){print((error+"\n"));
$trace8$;
exit(1);
$trace9$;
};
$trace10$;
returnres;
$trace11$;
};
$trace12$;
res=$open_f(file_path,m);
$trace13$;
returnres;
$trace14$;
};
$trace15$;
intread_file(intdescriptor_number,intsize,strings){intbytes;
$trace16$;
bytes=$read_f(descriptor_number,size,s);
$trace17$;
returnbytes;
$trace18$;
};
$trace19$;
voidclose_file(intdescriptor_number){bool$toPanicOld;
$trace20$;
$toPanicOld=$toPanic;
$trace21$;
$close_f(descriptor_number);
$trace22$;
try{strings;
$trace23$;
intbytes;
$trace24$;
bytes=read_file(descriptor_number,1,s);
$trace25$;
};
$trace26$;
string$s0;
$s0=str(descriptor_number);
if(error==""){error=("could not close file with file descriptor number "+$s0);
UNDEFINE($s0);
$trace27$;
if($toPanicOld){print((error+"\n"));
$trace28$;
exit(1);
$trace29$;
};
$trace30$;
};
$trace31$;
error="";
$trace32$;
};
$fileprogram.b$;
$trace1$;
voidmain(){intf;
$trace2$;
f=open_file("/home/slava/a1.txt","read");
$trace3$;
strings;
$trace4$;
intbytes;
$trace5$;
bytes=read_file(f,10,s);
$trace6$;
close_file(f);
$trace7$;
print((s+"\n"));
$trace8$;
};
$trace9$;
main();
