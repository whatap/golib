package fileutil

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
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

// ReadFile get file content
func ReadFile(filepath string, maxlength int64) ([]byte, int64, error) {
	f, e := os.Open(filepath)
	if e != nil {
		return nil, 0, e
	}
	defer f.Close()
	var output bytes.Buffer
	nbyteuntilnow := int64(0)

	// 버퍼를 동적으로 할당하여 필요한 만큼만 읽도록 수정
	for nbyteuntilnow < maxlength {
		// 남은 바이트와 기본 버퍼 크기 중 작은 값을 선택
		n := int64(4096)
		if n > maxlength-nbyteuntilnow {
			n = maxlength - nbyteuntilnow
		}
		buf := make([]byte, n)

		nbytethistime, e := f.Read(buf)
		if nbytethistime == 0 || e != nil {
			break
		}

		nbyteuntilnow += int64(nbytethistime)
		output.Write(buf[:nbytethistime])
	}

	if nbyteuntilnow > 0 {
		return output.Bytes(), nbyteuntilnow, nil
	}

	// 파일 내용이 0바이트인 경우 nil, 0, nil 반환
	// 에러가 없었으므로 e는 nil이 되어야 함
	return nil, 0, nil
}

// readEntireFile는 주어진 파일 경로의 내용을 모두 읽어 문자열로 반환합니다.
// 파일이 크더라도 메모리에 한 번에 모두 로드하여 처리합니다.
func ReadEntireFile(filePath string) (string, error) {
	// os.ReadFile을 사용하여 파일의 전체 내용을 바이트 슬라이스로 읽습니다.
	// 이 함수는 파일의 전체 내용을 한 번에 읽으므로, 파일이 매우 크다면 메모리 사용에 주의해야 합니다.
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("An error occurred while reading the file: %w", err)
	}

	// 읽어 들인 바이트 슬라이스를 문자열로 변환하여 반환합니다.
	return string(data), nil
}

// replaceFileWithOldBackup은 기존 파일을 새로운 내용으로 교체하고,
// 기존 파일은 old 파일로 백업합니다.
func ReplaceFileWithOldBackup(filePath, content string) error {
	// 1. 임시 파일 생성
	// filepath.Dir(filePath)를 사용하여 기존 파일이 있는 디렉토리에 임시 파일을 만듭니다.
	// .tmp 접미사를 사용하여 임시 파일임을 명확히 합니다.
	tmpFile, err := os.CreateTemp(filepath.Dir(filePath), "temp-*.tmp")
	if err != nil {
		return fmt.Errorf("Failed to create temporary file: %w", err)
	}
	// 함수 종료 시 임시 파일이 삭제되도록 defer를 사용합니다.
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// 2. 임시 파일에 새 내용 쓰기
	if _, err := io.WriteString(tmpFile, content); err != nil {
		return fmt.Errorf("Failed to write temporary file: %w", err)
	}

	// 3. 기존 파일이 이미 있는지 확인하고, 있다면 old 파일로 백업
	oldFilePath := filePath + ".old"
	if _, err := os.Stat(filePath); err == nil {
		// 파일이 존재하면 백업 파일(old)로 이름을 바꿉니다.
		// 만약 이미 old 파일이 있다면 덮어씁니다.
		if err := os.Rename(filePath, oldFilePath); err != nil {
			return fmt.Errorf("Failed to back up existing files: %w", err)
		}
		// 백업 파일이 더 이상 필요하지 않다면 os.Remove(oldFilePath)를 추가할 수 있습니다.
	}

	// 4. 임시 파일의 이름을 기존 파일 이름으로 변경
	if err := os.Rename(tmpFile.Name(), filePath); err != nil {
		return fmt.Errorf("Failed to rename temporary file: %w", err)
	}

	return nil
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
