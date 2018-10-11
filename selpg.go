package main

import (
	"io"
	"os"
	"os/exec"
	"fmt"
	"github.com/spf13/pflag"
	"strings"
	"math"
	"bufio"
)


type selpg_args struct {
	start_page int
	end_page int
	in_filename string
	page_len int
	page_type int
	print_dest string
}

var progname string

func usage() {
	fmt.Fprintf(os.Stderr,
		"\nUSAGE: %s -sstart_page -eend_page [ -f | -llines_per_page ] [ -ddest ] [ in_filename ]\n",
		 progname);
}

func process_args(ac int, psa *selpg_args) {
	if ac < 3 {
		fmt.Fprintf(os.Stderr, "%s: not enough arguments\n", progname)
		usage()
		os.Exit(1)
	}

	if strings.Contains(os.Args[1], "-s") != true {
		fmt.Fprintf(os.Stderr, "%s: 1st arg should be -sstart_page\n", progname)
		usage()
		os.Exit(2)
	}

	if psa.start_page < 1 || psa.start_page > math.MaxInt64-1 {
		fmt.Fprintf(os.Stderr, "%s: invalid start page %s\n", progname, string(psa.start_page))
		usage()
		os.Exit(3)		
	}

	if strings.Contains(os.Args[2], "-e") != true {
		fmt.Fprintf(os.Stderr, "%s: 2nd arg should be -eend_page\n", progname)
		usage()
		os.Exit(4)
	}

	if ( psa.end_page < 1 || psa.end_page > math.MaxInt64-1) || 
		(psa.end_page < psa.start_page) {
		fmt.Fprintf(os.Stderr, "%s: invalid end page %s\n", progname, string(psa.end_page))
		usage()
		os.Exit(5)
	}
}

func process_input(sa selpg_args) {
	fin := os.Stdin
	var fout io.Writer
	var line_ctr, page_ctr int
	var line string
	var c byte
	var err error
	var cmd *exec.Cmd
	if sa.in_filename != "" {
		fin, err = os.Open(sa.in_filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: could not open input file \"%s\"\n",
			progname, sa.in_filename)
			os.Exit(12)
		}
	}

	if sa.print_dest != "" {
		fout = bufio.NewWriter(os.Stdout)
	} else {
		fout = os.Stdout
	}

	input_reader := bufio.NewReader(fin)

	if sa.page_type == 'l' {
		line_ctr = 0;
		page_ctr = 1;

		for true {
			line, err = input_reader.ReadString('\n')
			if err != nil {
				break
			}
			line_ctr++
			if line_ctr > sa.page_len {
				page_ctr++
				line_ctr = 1
			}

			if (page_ctr >= sa.start_page) && page_ctr <= sa.end_page {
				fmt.Fprintf(fout, "%s", line)
			}

		}
	} else {
		page_ctr = 1
		for true {
			c, err = input_reader.ReadByte();
			if err != nil {
				break
			}
			if c == '\f' {
				page_ctr++;
			}
			if (page_ctr >= sa.start_page) && (page_ctr <= sa.end_page) {
				fmt.Fprintf(fout, "%c", c)
			}
		}
	}

	if page_ctr < sa.start_page {
		fmt.Fprintf(os.Stderr,
		"%s: start_page (%d) greater than total pages (%d), no output written\n",
		progname, sa.start_page, page_ctr);
	} else if page_ctr < sa.end_page {
		fmt.Fprintf(os.Stderr,
		"%s: end_page (%d) greater than total pages (%d), less output than expected\n",
		progname, sa.end_page, page_ctr);
	}
	
	if err == io.EOF {
		fin.Close()
		fmt.Fprintf(os.Stderr, "%s: done\n", progname)
	} else {
		fmt.Fprintf(os.Stderr, "%s: system error [%s] occurred on input stream fin\n",
		progname, err)
		fin.Close()
		os.Exit(14)
	}

	if sa.print_dest != "" {
		cmd = exec.Command("lp", "-d"+sa.print_dest)
		// cmd = exec.Command("cat", "-n")
		stdin, err := cmd.StdinPipe()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: could not open pipe to \"%s\"\n",
			progname, sa.print_dest)
		}

		go func() {
			defer stdin.Close()
			fmt.Fprint(stdin, fout)
		}()
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
		}
		fmt.Printf("%s\n", out)
	}
}

func main() {
	var sa selpg_args;
	pflag.IntVarP(&sa.start_page, "start", "s", 0, "start page")
	pflag.IntVarP(&sa.end_page, "end", "e", -1, "end page")
	pflag.IntVarP(&sa.page_len, "line", "l", 72, "line per page")
	pflag.IntVarP(&sa.page_type, "fomat", "f", 'l', "how to page")
	pflag.Lookup("fomat").NoOptDefVal = "1"
	pflag.StringVarP(&sa.print_dest, "dest", "d", "", "to printer")
	pflag.Parse();
	sa.in_filename = ""

	if pflag.NArg() != 0 {
		sa.in_filename = pflag.Arg(0)
	}

	progname = os.Args[0]

	process_args(len(os.Args), &sa)
	process_input(sa)

}
