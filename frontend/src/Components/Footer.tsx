import Link from "next/link";

const Footer = () => {
  return (
    <footer className="text-text-secondary text-center p-4 border-t-[1px] border-t-border text-sm">
      Made with ❤️ by{" "}
      <Link
        target="_blank"
        href="https://github.com/mehulzr"
        className="underline font-semibold"
      >
        Mehul
      </Link>
    </footer>
  );
};

export default Footer;
