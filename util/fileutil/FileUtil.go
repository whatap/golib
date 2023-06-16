package fileutil

import (
	"bytes"
	"fmt"
	"net"
	"os"
)

func CloseOfDatagramSocket(udp *net.UDPConn) error {
	if udp != nil {
		return udp.Close()
	}
	return fmt.Errorf("%s", "UDPConn is nil")
}

func ExistsFile(filename string) (os.FileInfo, bool) {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return nil, false
	}
	if err != nil {
		return nil, false
	}
	return info, !info.IsDir()
}

func ExistsDir(filename string) (os.FileInfo, bool) {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return nil, false
	}
	if err != nil {
		return nil, false
	}
	return info, info.IsDir()
}
func Info(filename string) (os.FileInfo, error) {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return nil, err
	}
	return info, nil
}

//ReadFile get file content
func ReadFile(filepath string, maxlength int64) ([]byte, int64, error) {
	f, e := os.Open(filepath)
	if e != nil {
		return nil, 0, e
	}
	defer f.Close()
	var output bytes.Buffer
	buf := make([]byte, 4096)
	nbyteuntilnow := int64(0)
	for nbyteleft := maxlength; nbyteleft > 0; {
		nbytethistime, e := f.Read(buf)
		if nbytethistime == 0 || e != nil {
			break
		}
		nbyteleft -= int64(nbytethistime)
		nbyteuntilnow += int64(nbytethistime)
		output.Write(buf[:nbytethistime])
	}

	if nbyteuntilnow > 0 {
		return output.Bytes(), nbyteuntilnow, nil
	}

	return nil, 0, e
}

//
//	public static InputStream close(InputStream in) {
//		try {
//			if (in != null) {
//				in.close();
//			}
//		} catch (Throwable e) {
//		}
//		return null;
//	}
//
//	public static OutputStream close(OutputStream out) {
//		try {
//			if (out != null) {
//				out.close();
//			}
//		} catch (Throwable e) {
//		}
//		return null;
//	}
//
//	public static Reader close(Reader in) {
//		try {
//			if (in != null) {
//				in.close();
//			}
//		} catch (Throwable e) {
//		}
//		return null;
//	}
//
//	public static Writer close(Writer out) {
//		try {
//			if (out != null) {
//				out.close();
//			}
//		} catch (Throwable e) {
//		}
//		return null;
//	}
//
//	public static byte[] readAll(InputStream fin) throws IOException {
//		ByteArrayOutputStream out = new ByteArrayOutputStream();
//		byte[] buff = new byte[4096];
//		int n = fin.read(buff);
//		while (n >= 0) {
//			out.write(buff, 0, n);
//			n = fin.read(buff);
//		}
//		return out.toByteArray();
//	}
//
//	public static IClose close(IClose object) {
//		try {
//			if (object != null) {
//				object.close();
//			}
//		} catch (Throwable e) {
//			e.printStackTrace();
//		}
//		return null;
//	}
//
//	public static RandomAccessFile close(RandomAccessFile raf) {
//		try {
//			if (raf != null) {
//				raf.close();
//			}
//		} catch (Throwable e) {
//		}
//		return null;
//	}
//
//	public static Socket close(Socket socket) {
//		try {
//			if (socket != null) {
//				socket.close();
//			}
//		} catch (Throwable e) {
//		}
//		return null;
//	}
//
//	public static ServerSocket close(ServerSocket socket) {
//		try {
//			if (socket != null) {
//				socket.close();
//			}
//		} catch (Throwable e) {
//		}
//		return null;
//	}
//
//	public static void save(String file, byte[] b) {
//		save(new File(file), b);
//	}
//
//	public static void save(File file, byte[] byteArray) {
//		FileOutputStream out = null;
//		try {
//			out = new FileOutputStream(file);
//			out.write(byteArray);
//		} catch (Exception e) {
//		}
//		close(out);
//	}
//
//	public static byte[] readAll(File file) {
//		FileInputStream in = null;
//		try {
//			in = new FileInputStream(file);
//			return readAll(in);
//		} catch (Exception e) {
//		} finally {
//			close(in);
//		}
//		return null;
//	}
//
//	public static void copy(File src, File dest) {
//		try {
//			copy(src, dest, true);
//		} catch (Exception e) {
//		}
//	}
//
//	public static boolean copy(File src, File dest, boolean overwrite) throws IOException {
//		if (!src.isFile() || !src.exists())
//			return false;
//		if (dest.exists()) {
//			if (dest.isDirectory()) // Directory이면 src파일명을 사용한다.
//				dest = new File(dest, src.getName());
//			else if (dest.isFile()) {
//				if (!overwrite)
//					throw new IOException(dest.getAbsolutePath() + "' already exists!");
//			} else
//				throw new IOException("Invalid  file '" + dest.getAbsolutePath() + "'");
//		}
//		File destDir = dest.getParentFile();
//		if (!destDir.exists())
//			if (!destDir.mkdirs())
//				throw new IOException("Failed to create " + destDir.getAbsolutePath());
//		long fileSize = src.length();
//		if (fileSize > 20 * 1024 * 1024) {
//			FileInputStream in = null;
//			FileOutputStream out = null;
//			try {
//				in = new FileInputStream(src);
//				out = new FileOutputStream(dest);
//				int done = 0;
//				int buffLen = 32768;
//				byte buf[] = new byte[buffLen];
//				while ((done = in.read(buf, 0, buffLen)) >= 0) {
//					if (done == 0)
//						Thread.yield();
//					else
//						out.write(buf, 0, done);
//				}
//			} finally {
//				close(in);
//				close(out);
//			}
//		} else {
//			FileInputStream in = null;
//			FileOutputStream out = null;
//			FileChannel fin = null;
//			FileChannel fout = null;
//			try {
//				in = new FileInputStream(src);
//				out = new FileOutputStream(dest);
//				fin = in.getChannel();
//				fout = out.getChannel();
//				long position = 0;
//				long done = 0;
//				long count = Math.min(65536, fileSize);
//				do {
//					done = fin.transferTo(position, count, fout);
//					position += done;
//					fileSize -= done;
//				} while (fileSize > 0);
//			} finally {
//				close(fin);
//				close(fout);
//				close(in);
//				close(out);
//			}
//		}
//		return true;
//	}
//
//	public static FileChannel close(FileChannel fc) {
//		if (fc != null) {
//			try {
//				fc.close();
//			} catch (IOException e) {
//			}
//		}
//		return null;
//	}
//
//	// public static void chmod777(File file) {
//	// try {
//	// file.setReadable(true, false);
//	// file.setWritable(true, false);
//	// file.setExecutable(true, false);
//	// } catch (Throwable th) {}
//	// }
//	public static void close(DataInputX in) {
//		try {
//			if (in != null)
//				in.close();
//		} catch (Exception e) {
//		}
//	}
//
//	public static void close(DataOutputX out) {
//		try {
//			if (out != null)
//				out.close();
//		} catch (Exception e) {
//		}
//	}
//
//	public static String load(File file, String enc) {
//		if (file == null || file.canRead() == false)
//			return null;
//		BufferedInputStream in = null;
//		try {
//			in = new BufferedInputStream(new FileInputStream(file));
//			return new String(readAll(in), enc);
//		} catch (IOException e) {
//			e.printStackTrace();
//		} finally {
//			close(in);
//		}
//		return null;
//	}
//
//	public static void append(String file, String line) {
//		PrintWriter out = null;
//		try {
//			out = new PrintWriter(new FileWriter(file, true));
//			out.println(line);
//		} catch (Exception e) {
//		}
//		close(out);
//	}
//
//	public static boolean mkdirs(String path) {
//		File f = new File(path);
//		if (f.exists() == false)
//			return f.mkdirs();
//		else
//			return true;
//	}
//
//	public static Properties readProperties(File f) {
//		BufferedInputStream reader = null;
//		Properties p = new Properties();
//		try {
//			reader = new BufferedInputStream(new FileInputStream(f));
//			p.load(reader);
//		} catch (Exception e) {
//		} finally {
//			close(reader);
//		}
//		return p;
//	}
//
//	public static void writeProperties(File f, Properties p) {
//		PrintWriter pw = null;
//		try {
//			pw = new PrintWriter(f);
//			p.list(pw);
//		} catch (Exception e) {
//		} finally {
//			close(pw);
//		}
//	}
//
//	public static void close(Connection conn) {
//		try {
//			if (conn != null) {
//				conn.close();
//			}
//		} catch (Throwable e) {
//			e.printStackTrace();
//		}
//	}
//
//	public static void close(Statement stmt) {
//		try {
//			if (stmt != null) {
//				stmt.close();
//			}
//		} catch (Throwable e) {
//			e.printStackTrace();
//		}
//	}
//
//	public static void close(ResultSet rs) {
//		try {
//			if (rs != null) {
//				rs.close();
//			}
//		} catch (Throwable e) {
//			e.printStackTrace();
//		}
//	}
//
//	public static void close(JarFile jarfile) {
//		try {
//			if (jarfile != null) {
//				jarfile.close();
//			}
//		} catch (Throwable e) {
//		}
//	}
//
//	public static void delete(String home, String prefix, String postfix) {
//		File[] files = new File(home).listFiles();
//		for (int i = 0; i < files.length; i++) {
//			if (files[i].isFile()) {
//				String nm = files[i].getName();
//				if (nm.startsWith(prefix) && nm.endsWith(postfix)) {
//					files[i].delete();
//				}
//			}
//		}
//	}
//
//	public static void createNew(File file, long size) {
//		BufferedOutputStream w = null;
//		try {
//			w = new BufferedOutputStream(new FileOutputStream(file));
//			for (int i = 0; i < size; i++) {
//				w.write(0);
//			}
//		} catch (Exception e) {
//
//		} finally {
//			close(w);
//		}
//	}
//
//	public static int getCrc(File f) {
//		BufferedInputStream in = null;
//		try {
//			in = new BufferedInputStream(new FileInputStream(f));
//			return HashUtil.hash(in, (int) f.length());
//		} catch (Exception e) {
//		} finally {
//			close(in);
//		}
//		return 0;
//	}
